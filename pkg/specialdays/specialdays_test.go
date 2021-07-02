package specialdays

import (
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	testCases := []struct {
		Description   string
		Contacts      []map[string]interface{}
		CheckDate     time.Time
		Alert         bool
		ExpectedTitle string
		ExpectedBody  string
	}{
		{
			Description: "it sets nothing when there are no days",
			Contacts: []map[string]interface{}{
				{
					"Display Name": "John Appleseed",
				},
			},
			Alert:         false,
			ExpectedTitle: "",
			ExpectedBody:  "",
		},
		{
			Description: "it alerts on birthdays",
			Contacts: []map[string]interface{}{
				{
					"Display Name": "John Appleseed",
					"Birthday":     "1995-06-03",
				},
			},
			CheckDate:     time.Date(2021, time.June, 3, 0, 0, 0, 0, time.UTC),
			Alert:         true,
			ExpectedTitle: "John Appleseed's Birthday",
			ExpectedBody:  "It's John Appleseed's birthday",
		},
		{
			Description: "it alerts on multiple birthdays",
			Contacts: []map[string]interface{}{
				{
					"Display Name": "John",
					"Birthday":     "1995-06-03",
				},
				{
					"Display Name": "Jane",
					"Birthday":     "1995-06-03",
				},
				{
					"Display Name": "Jill",
					"Birthday":     "1995-06-03",
				},
			},
			CheckDate:     time.Date(2021, time.June, 3, 0, 0, 0, 0, time.UTC),
			Alert:         true,
			ExpectedTitle: "3 birthdays",
			ExpectedBody:  "John, Jane & Jill have birthdays",
		},
		{
			Description: "it alerts on special days",
			Contacts: []map[string]interface{}{
				{
					"Display Name":      "John",
					"JSON Special Days": `[{ "date": "1994-03-21", "label": "aniversary"},{ "date": "1994-06-02", "label": "other"}]`,
				},
			},
			CheckDate:     time.Date(2021, time.March, 21, 0, 0, 0, 0, time.UTC),
			Alert:         true,
			ExpectedTitle: "John has a special day labelled 'aniversary'",
			ExpectedBody:  "John has a special day labelled 'aniversary'",
		},
		{
			Description: "it alerts on multiple special days",
			Contacts: []map[string]interface{}{
				{
					"Display Name":      "John",
					"JSON Special Days": `[{ "date": "1994-03-21", "label": "aniversary"},{ "date": "1994-06-02", "label": "other"}]`,
				},
				{
					"Display Name":      "Jane",
					"JSON Special Days": `[{ "date": "1994-03-21", "label": "other"}]`,
				},
			},
			CheckDate:     time.Date(2021, time.March, 21, 0, 0, 0, 0, time.UTC),
			Alert:         true,
			ExpectedTitle: "John & Jane have special days",
			ExpectedBody: `* John has a special day labelled 'aniversary'
* Jane has a special day labelled 'other'`,
		},
		{
			Description: "it alerts on birthdays and special days",
			Contacts: []map[string]interface{}{
				{
					"Display Name":      "John",
					"JSON Special Days": `[{ "date": "1994-03-21", "label": "aniversary"},{ "date": "1994-06-02", "label": "other"}]`,
				},
				{
					"Display Name":      "Jane",
					"JSON Special Days": `[{ "date": "1994-03-21", "label": "other"}]`,
				},
				{
					"Display Name": "Jill",
					"Birthday":     "1995-03-21",
				},
			},
			CheckDate:     time.Date(2021, time.March, 21, 0, 0, 0, 0, time.UTC),
			Alert:         true,
			ExpectedTitle: "John, Jane & Jill have events",
			ExpectedBody: `* John has a special day labelled 'aniversary'
* Jane has a special day labelled 'other'
* Jill has a birthday`,
		},
	}

	for _, test := range testCases {
		t.Run(test.Description, func(t *testing.T) {
			alert, title, body, err := Generate(test.Contacts, test.CheckDate)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if test.Alert != alert {
				t.Fatalf("unexpected alert:\ngot %v want %v", alert, test.Alert)
			}
			if test.ExpectedTitle != title {
				t.Errorf("unexpected title:\ngot %v want %v", title, test.ExpectedTitle)
			}
			if test.ExpectedBody != body {
				t.Errorf("unexpected body:\ngot\n%v\nwant\n%v", body, test.ExpectedBody)
			}
		})
	}
}
