package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/charlieegan3/tool-airtable-contacts/pkg/airtable"
	"github.com/charlieegan3/tool-airtable-contacts/pkg/carddav"
	air "github.com/mehanizm/airtable"
)

// Sync copies the data in airtable into a carddav server for sycning to devices
type Sync struct {
	ScheduleOverride string
	Endpoint         string

	AirtableKey   string
	AirtableBase  string
	AirtableTable string
	AirtableView  string

	CardDAVServer   string
	CardDAVUser     string
	CardDAVPassword string

	VCardPhotoSize int
	VCardV3        bool
}

func (s *Sync) Name() string {
	return "sync"
}

func (s *Sync) Run(ctx context.Context) error {
	doneCh := make(chan bool)
	errCh := make(chan error)

	go func() {
		// get the latest data
		airtableClient := air.NewClient(s.AirtableKey)
		records, err := airtable.Download(airtableClient, s.AirtableBase, s.AirtableTable, s.AirtableView)
		if err != nil {
			errCh <- fmt.Errorf("failed to download contacts: %s", err)
			return
		}

		cardDAVClient := carddav.Client{
			URL:      s.CardDAVServer,
			User:     s.CardDAVUser,
			Password: s.CardDAVPassword,
		}

		// records passed as vcard sync is done on per contact basis
		err = carddav.Sync(cardDAVClient, records, s.VCardV3, s.VCardPhotoSize)
		if err != nil {
			errCh <- fmt.Errorf("failed to upload to carddav: %s", err)
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

func (s *Sync) Timeout() time.Duration {
	return 60 * time.Second
}

func (s *Sync) Schedule() string {
	if s.ScheduleOverride != "" {
		return s.ScheduleOverride
	}
	return "0 */15 * * * 0"
}
