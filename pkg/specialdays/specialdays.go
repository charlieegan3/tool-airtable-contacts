package specialdays

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func Generate(contacts []map[string]interface{}, checkDate time.Time, period int) (bool, string, string, error) {
	// likely invalid call, just return nothing
	if period == 0 {
		return false, "", "", nil
	}

	periodStart := checkDate
	periodEnd := checkDate.Add(time.Duration(period) * time.Hour * 24).Add(-1 * time.Second)

	contactsWithBirthdays := []string{}
	contactsWithSpecialDays := []string{}
	specialDayMessages := []string{}

	for _, contact := range contacts {
		// list names with birthdays
		birthdayValue, ok := contact["Birthday"].(string)
		if ok {
			birthday, err := time.Parse("2006-01-02", birthdayValue)
			if err != nil {
				return false, "", "", fmt.Errorf("failed to parse birthday value for: %v", contact)
			}

			birthdayThisYear := time.Date(time.Now().UTC().Year(), birthday.Month(), birthday.Day(), 0, 0, 0, 0, time.UTC)

			if dateInPeriod(periodStart, periodEnd, birthdayThisYear) {
				contactsWithBirthdays = append(contactsWithBirthdays, contact["Display Name"].(string))
			}
		}

		// list names for special days
		jsonSpecialDaysValue, ok := contact["JSON Special Days"].(string)
		if ok {
			decoder := json.NewDecoder(strings.NewReader(jsonSpecialDaysValue))
			days := []struct {
				Label string `json:"label"`
				Date  string `json:"date"`
			}{}
			err := decoder.Decode(&days)
			if err != nil {
				return false, "", "", fmt.Errorf("failed to parse JSON Special Days value for: %v", contact)
			}
			for _, day := range days {
				date, err := time.Parse("2006-01-02", day.Date)
				if err != nil {
					return false, "", "", fmt.Errorf("failed to parse special day date value for: %v", contact)
				}

				dateThisYear := time.Date(time.Now().UTC().Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

				if dateInPeriod(periodStart, periodEnd, dateThisYear) {
					contactsWithSpecialDays = append(contactsWithSpecialDays, contact["Display Name"].(string))
					specialDayMessages = append(
						specialDayMessages,
						fmt.Sprintf("%s has a special day labelled '%s'", contact["Display Name"].(string), day.Label),
					)
				}
			}
		}
	}

	if len(contactsWithSpecialDays) > 0 && len(contactsWithBirthdays) > 0 {
		allNames := append(contactsWithSpecialDays, contactsWithBirthdays...)
		allMessages := []string{}

		for _, message := range specialDayMessages {
			allMessages = append(allMessages, fmt.Sprintf("* %s", message))
		}
		for _, name := range contactsWithBirthdays {
			allMessages = append(allMessages, fmt.Sprintf("* %s has a birthday", name))
		}

		initialList := strings.Join(allNames[0:len(allNames)-1], ", ")
		return true,
			fmt.Sprintf("%s & %s have events", initialList, allNames[len(allNames)-1]),
			strings.Join(allMessages, "\n"),
			nil
	}

	if len(specialDayMessages) == 1 {
		return true,
			specialDayMessages[0],
			specialDayMessages[0],
			nil
	}

	if len(specialDayMessages) > 1 {
		initialList := strings.Join(contactsWithSpecialDays[0:len(contactsWithSpecialDays)-1], ", ")
		return true,
			fmt.Sprintf("%s & %s have special days", initialList, contactsWithSpecialDays[len(contactsWithSpecialDays)-1]),
			fmt.Sprintf("* %s", strings.Join(specialDayMessages, "\n* ")),
			nil
	}

	if len(contactsWithBirthdays) == 1 {
		return true,
			fmt.Sprintf("%s's Birthday", contactsWithBirthdays[0]),
			fmt.Sprintf("It's %s's birthday", contactsWithBirthdays[0]),
			nil
	}

	if len(contactsWithBirthdays) > 1 {
		initialList := strings.Join(contactsWithBirthdays[0:len(contactsWithBirthdays)-1], ", ")
		return true,
			fmt.Sprintf("%d birthdays", len(contactsWithBirthdays)),
			fmt.Sprintf("%s & %s have birthdays", initialList, contactsWithBirthdays[len(contactsWithBirthdays)-1]),
			nil
	}

	return false, "", "", nil
}

func dateInPeriod(periodStart, periodEnd, date time.Time) bool {
	if date.Equal(periodStart) || date.Equal(periodEnd) {
		return true
	}
	if date.After(periodStart) && date.Before(periodEnd) {
		return true
	}
	return false
}
