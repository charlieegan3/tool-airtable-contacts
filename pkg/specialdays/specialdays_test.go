package specialdays

import (
	"fmt"
	"testing"
	"time"
)

func TestDateInPeriod(t *testing.T) {
	testCases := []struct {
		PeriodStart time.Time
		PeriodEnd   time.Time
		Date        time.Time
		Result      bool
	}{
		{
			PeriodStart: time.Date(2021, time.March, 21, 0, 0, 0, 0, time.UTC),
			PeriodEnd:   time.Date(2021, time.March, 23, 0, 0, 0, 0, time.UTC),
			Date:        time.Date(2021, time.March, 22, 0, 0, 0, 0, time.UTC),
			Result:      true,
		},
		{
			PeriodStart: time.Date(2021, time.March, 21, 0, 0, 0, 0, time.UTC),
			PeriodEnd:   time.Date(2021, time.March, 21, 0, 0, 0, 0, time.UTC),
			Date:        time.Date(2021, time.March, 21, 0, 0, 0, 0, time.UTC),
			Result:      true,
		},
		{
			PeriodStart: time.Date(2021, time.March, 21, 0, 0, 0, 0, time.UTC),
			PeriodEnd:   time.Date(2021, time.March, 23, 0, 0, 0, 0, time.UTC),
			Date:        time.Date(2021, time.March, 23, 0, 0, 0, 0, time.UTC),
			Result:      true,
		},
		{
			PeriodStart: time.Date(2021, time.March, 21, 0, 0, 0, 0, time.UTC),
			PeriodEnd:   time.Date(2021, time.March, 23, 0, 0, 0, 0, time.UTC),
			Date:        time.Date(2021, time.March, 24, 0, 0, 0, 0, time.UTC),
			Result:      false,
		},
		{
			PeriodStart: time.Date(2021, time.March, 21, 0, 0, 0, 0, time.UTC),
			PeriodEnd:   time.Date(2021, time.March, 23, 0, 0, 0, 0, time.UTC),
			Date:        time.Date(2021, time.March, 20, 0, 0, 0, 0, time.UTC),
			Result:      false,
		},
	}

	for i, test := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result := dateInPeriod(test.PeriodStart, test.PeriodEnd, test.Date)
			if result != test.Result {
				t.Fatalf("unexpected result: %v", result)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	testCases := []struct {
		Description   string
		Contacts      []map[string]interface{}
		CheckDate     time.Time
		Period        int
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
			Period:        1,
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
				{
					"Display Name": "Nigel",
					"Birthday":     "1995-06-04",
				},
			},
			CheckDate:     time.Date(2021, time.June, 3, 0, 0, 0, 0, time.UTC),
			Period:        1,
			Alert:         true,
			ExpectedTitle: "3 birthdays",
			ExpectedBody:  "John, Jane & Jill have birthdays",
		},
		{
			Description: "it alerts on birthdays in period",
			Contacts: []map[string]interface{}{
				{
					"Display Name": "John",
					"Birthday":     "1995-06-01",
				},
				{
					"Display Name": "Jane",
					"Birthday":     "1995-06-04",
				},
				{
					"Display Name": "Jill",
					"Birthday":     "1995-06-05",
				},
			},
			CheckDate:     time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC),
			Period:        5,
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
			Period:        1,
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
			Period:        1,
			Alert:         true,
			ExpectedTitle: "John & Jane have special days",
			ExpectedBody: `* John has a special day labelled 'aniversary'
* Jane has a special day labelled 'other'`,
		},
		{
			Description: "it alerts on special days in period",
			Contacts: []map[string]interface{}{
				{
					"Display Name":      "John",
					"JSON Special Days": `[{ "date": "1994-03-10", "label": "aniversary"},{ "date": "1994-06-02", "label": "other"}]`,
				},
				{
					"Display Name":      "Jane",
					"JSON Special Days": `[{ "date": "1994-03-01", "label": "other"}]`,
				},
			},
			CheckDate:     time.Date(2021, time.March, 1, 0, 0, 0, 0, time.UTC),
			Period:        10,
			Alert:         true,
			ExpectedTitle: "John & Jane have special days",
			ExpectedBody: `* John has a special day labelled 'aniversary'
* Jane has a special day labelled 'other'`,
		},
		{
			Description: "it alerts on named special days",
			Contacts: []map[string]interface{}{
				{
					"Display Name":      "John",
					"JSON Special Days": `[{ "date": "fathers-day-uk", "label": "fathers-day-uk"}]`,
				},
				{
					"Display Name":      "Jane",
					"JSON Special Days": `[{ "date": "mothering-sunday", "label": "mothering-sunday"}]`,
				},
			},
			CheckDate:     time.Date(2021, time.March, 1, 0, 0, 0, 0, time.UTC),
			Period:        200, // just a big period to catch and test the above
			Alert:         true,
			ExpectedTitle: "John & Jane have special days",
			ExpectedBody: `* John has a special day labelled 'fathers-day-uk'
* Jane has a special day labelled 'mothering-sunday'`,
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
			Period:        1,
			Alert:         true,
			ExpectedTitle: "John, Jane & Jill have events",
			ExpectedBody: `* John has a special day labelled 'aniversary'
* Jane has a special day labelled 'other'
* Jill has a birthday`,
		},
		{
			Description: "it alerts on birthdays and special days in period",
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
			CheckDate:     time.Date(2021, time.March, 1, 0, 0, 0, 0, time.UTC),
			Period:        30,
			Alert:         true,
			ExpectedTitle: "John, Jane & Jill have events",
			ExpectedBody: `* John has a special day labelled 'aniversary'
* Jane has a special day labelled 'other'
* Jill has a birthday`,
		},
	}

	for _, test := range testCases {
		t.Run(test.Description, func(t *testing.T) {
			alert, title, body, err := Generate(test.Contacts, test.CheckDate, test.Period)
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
