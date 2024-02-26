package airtable

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	air "github.com/mehanizm/airtable"
)

func Download(client *air.Client, databaseID, tableName, viewName string) ([]map[string]interface{}, error) {
	var records []map[string]interface{}

	table := client.GetTable(databaseID, tableName)

	offset := "0"
	for {
		var result *air.Records
		var err error
		operation := func() error {
			result, err = table.GetRecords().
				FromView(viewName).
				WithOffset(offset).
				ReturnFields("ID", "Display Name", "JSON Addresses", "JSON Phone Numbers", "JSON Emails", "Note", "Company", "Profile Image", "Birthday", "JSON Special Days").
				Do()
			if err != nil {
				return fmt.Errorf("failed to get records: %s", err)
			}

			for _, v := range result.Records {
				records = append(records, v.Fields)
			}

			return nil
		}

		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = 3 * time.Minute

		err = backoff.Retry(operation, b)
		if err != nil {
			return records, fmt.Errorf("failed to get records after backoff: %s", err)
		}

		// have reached the end
		if result.Offset == "" {
			break
		}

		offset = result.Offset
	}

	return records, nil
}
