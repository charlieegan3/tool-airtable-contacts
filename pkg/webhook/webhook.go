package webhook

import (
	"bytes"
	json2 "encoding/json"
	"fmt"
	"net/http"
)

func Send(endpoint, title, body, linkURL string) error {
	data := []map[string]string{
		{
			"Title": title,
			"Body":  body,
			"URL":   linkURL,
		},
	}

	b, err := json2.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook body: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}

	return nil
}
