package carddav

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
)

// Client wraps functionality for calling carddav methods
type Client struct {
	URL      string
	User     string
	Password string
}

func (c *Client) List() (items []string, err error) {
	body, err := c.do(
		c.URL,
		"PROPFIND",
		map[string]string{
			"Depth": "1",
		},
		nil,
	)
	if err != nil {
		return items, errors.Wrap(err, "failed to list carddav")
	}
	defer body.Close()

	doc, err := xmlquery.Parse(body)
	if err != nil {
		return items, errors.Wrap(err, "failed parse body as xml")
	}
	list, err := xmlquery.QueryAll(doc, "D:multistatus/D:response/D:href")
	if err != nil {
		return items, errors.Wrap(err, "failed to query body for hrefs")
	}
	for _, v := range list {
		// only take vcf files from the list
		if strings.HasSuffix(v.InnerText(), "vcf") {
			id, err := extractID(v.InnerText())
			if err != nil {
				return items, errors.Wrap(err, "failed to get response item ID")
			}
			items = append(items, id)
		}
	}

	return items, nil
}

func (c *Client) Delete(id string) (err error) {
	body, err := c.do(
		fmt.Sprintf("%s/%s.vcf", strings.TrimSuffix(c.URL, "/"), id),
		"DELETE",
		map[string]string{
			"Depth": "1",
		},
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to delete vcf")
	}

	defer body.Close()

	return nil
}

func (c *Client) Put(id string, vcardData string) (err error) {
	body, err := c.do(
		fmt.Sprintf("%s/%s.vcf", strings.TrimSuffix(c.URL, "/"), id),
		"PUT",
		map[string]string{
			"Content-Type": "text/vcard",
		},
		strings.NewReader(vcardData),
	)
	if err != nil {
		return errors.Wrap(err, "failed to put vcf file")
	}
	defer body.Close()

	return nil
}

func (c *Client) do(url string, verb string, headers map[string]string, requestBody io.Reader) (io.ReadCloser, error) {
	var body io.ReadCloser
	operation := func() error {
		req, err := http.NewRequest(verb, url, requestBody)
		if err != nil {
			return errors.Wrap(err, "failed to make PROPFIND request to list items in carddav endpoint")
		}

		req.SetBasicAuth(c.User, c.Password)

		for k, v := range headers {
			req.Header.Add(k, v)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return errors.Wrap(err, "failed to make request")
		}

		if resp.StatusCode > 399 || resp.StatusCode < 100 {
			return fmt.Errorf("server returned error: %d", resp.StatusCode)
		}

		body = resp.Body

		return nil
	}

	b := backoff.NewExponentialBackOff()

	err := backoff.Retry(operation, b)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list carddav items after backoff")
	}

	return body, nil
}

func extractID(fullPath string) (id string, err error) {
	parts := strings.Split(fullPath, "/")
	if len(parts) != 7 {
		return id, fmt.Errorf("unexpected number of path items: %d", len(parts))
	}

	return strings.TrimSuffix(parts[6], ".vcf"), nil
}
