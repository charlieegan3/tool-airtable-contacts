package carddav

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
)

// Client wraps functionality for calling carddav methods
type Client struct {
	URL      string
	User     string
	Password string
}

func (c *Client) List() (items []string, err error) {
	req, err := http.NewRequest("PROPFIND", c.URL, nil)
	if err != nil {
		return items, errors.Wrap(err, "failed to make PROPFIND request to list items in carddav endpoint")
	}
	req.SetBasicAuth(c.User, c.Password)
	req.Header.Add("Depth", "1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return items, errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 || resp.StatusCode < 100 {
		return items, errors.Wrap(err, "server returned error: ")
	}

	doc, err := xmlquery.Parse(resp.Body)
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
	url := fmt.Sprintf("%s/%s.vcf", strings.TrimSuffix(c.URL, "/"), id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to make DELETE request")
	}
	req.SetBasicAuth(c.User, c.Password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 || resp.StatusCode < 100 {
		return errors.Wrap(err, "server returned error: ")
	}
	return nil
}

func (c *Client) Put(id string, vcardData string) (err error) {
	url := fmt.Sprintf("%s/%s.vcf", strings.TrimSuffix(c.URL, "/"), id)

	req, err := http.NewRequest("PUT", url, strings.NewReader(vcardData))
	if err != nil {
		return errors.Wrap(err, "failed to make PUT request")
	}
	req.SetBasicAuth(c.User, c.Password)
	req.Header.Add("Content-Type", "text/vcard")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 || resp.StatusCode < 100 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read error body")
		}
		return fmt.Errorf("server returned error: %d, %s", resp.StatusCode, body)
	}

	return nil
}

func extractID(fullPath string) (id string, err error) {
	parts := strings.Split(fullPath, "/")
	if len(parts) != 7 {
		return id, fmt.Errorf("unexpected number of path items: %d", len(parts))
	}

	return strings.TrimSuffix(parts[6], ".vcf"), nil
}
