package airtable

import (
	"fmt"

	air "github.com/mehanizm/airtable"
)

func Download(client *air.Client, databaseID, tableName, viewName string) ([]map[string]interface{}, error) {
	var records []map[string]interface{}

	table := client.GetTable(databaseID, tableName)

	offset := ""
	for {
		result, err := table.GetRecords().
			FromView(viewName).
			WithOffset(offset).
			ReturnFields("Display Name", "Emails", "JSON Phone Numbers", "Note", "Company", "Profile Image").
			Do()
		if err != nil {
			return records, fmt.Errorf("failed to get records: %s", err)
		}

		for _, v := range result.Records {
			records = append(records, v.Fields)
		}

		// have reached the end
		if len(result.Records) < 100 {
			break
		}

		offset = result.Offset
	}

	return records, nil
}
