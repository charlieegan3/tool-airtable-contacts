package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/charlieegan3/tool-airtable-contacts/pkg/specialdays"

	"github.com/charlieegan3/tool-airtable-contacts/pkg/airtable"
	"github.com/gomarkdown/markdown"
	air "github.com/mehanizm/airtable"
)

// Day sends any alerts for the current day
type Day struct {
	ScheduleOverride string
	Endpoint         string

	AirtableKey   string
	AirtableBase  string
	AirtableTable string
	AirtableView  string
}

func (d *Day) Name() string {
	return "notify-day"
}

func (d *Day) Run(ctx context.Context) error {
	doneCh := make(chan bool)
	errCh := make(chan error)

	go func() {
		// get the latest data
		airtableClient := air.NewClient(d.AirtableKey)
		records, err := airtable.Download(airtableClient, d.AirtableBase, d.AirtableTable, d.AirtableView)
		if err != nil {
			errCh <- fmt.Errorf("failed to download contacts: %s", err)
			return
		}

		// set the notification period
		periodStart := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
		alert, title, body, err := specialdays.Generate(records, periodStart, 1, true)
		if err != nil {
			errCh <- fmt.Errorf("failed to generate alert message: %s", err)
			return
		}

		// send the alert if needed
		if !alert {
			doneCh <- true
			return
		}

		datab := []map[string]string{
			{
				"title": fmt.Sprintf("Today: %s", title),
				"body":  string(markdown.ToHTML([]byte(body), nil, nil)),
				"url":   "",
			},
		}

		b, err := json.Marshal(datab)
		if err != nil {
			errCh <- fmt.Errorf("failed to form new item JSON: %s", err)
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", d.Endpoint, bytes.NewBuffer(b))
		if err != nil {
			errCh <- fmt.Errorf("failed to build request for new item: %s", err)
			return
		}

		req.Header.Add("Content-Type", "application/json; charset=utf-8")

		resp, err := client.Do(req)
		if err != nil {
			errCh <- fmt.Errorf("failed to send request for new item: %s", err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			errCh <- fmt.Errorf("failed to send request: non 200OK response")
			return
		}

		doneCh <- true
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-errCh:
		return fmt.Errorf("job failed with error: %s", e)
	case <-doneCh:
		return nil
	}
}

func (d *Day) Timeout() time.Duration {
	return 30 * time.Second
}

func (d *Day) Schedule() string {
	if d.ScheduleOverride != "" {
		return d.ScheduleOverride
	}
	return "0 0 5 * * *"
}
