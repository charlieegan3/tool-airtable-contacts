package specialdays

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/charlieegan3/special-days/pkg/fathersday"
	"github.com/charlieegan3/special-days/pkg/motheringsunday"
)

func Generate(contacts []map[string]interface{}, checkDate time.Time, period int, today bool) (bool, string, string, error) {
	// likely invalid call, just return nothing
	if period == 0 {
		return false, "", "", nil
	}

	periodStart := checkDate
	periodEnd := checkDate.Add(time.Duration(period) * time.Hour * 24).Add(-1 * time.Second)

	contactsWithBirthdays := []string{}
	birthdayMessages := []string{}
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

			// calculate the birthday date for the year of the period start
			birthdayInYearForPeriod := time.Date(
				periodStart.Year(),
				birthday.Month(),
				birthday.Day(),
				0, 0, 0, 0,
				time.UTC,
			)

			if dateInPeriod(periodStart, periodEnd, birthdayInYearForPeriod) {
				name := contact["Display Name"].(string)
				contactsWithBirthdays = append(contactsWithBirthdays, name)
				message := fmt.Sprintf("%s has a birthday", name)
				if !today {
					message = fmt.Sprintf("%s on %s", message, birthday.Format("January 02"))
				}
				birthdayMessages = append(birthdayMessages, message)
			}
		}

		// list names for special days
		jsonSpecialDaysValue, ok := contact["JSON Special Days"].(string)
		if ok {
			decoder := json.NewDecoder(strings.NewReader(jsonSpecialDaysValue))
			var days []struct {
				Label string `json:"label"`
				Date  string `json:"date"`
			}
			err := decoder.Decode(&days)
			if err != nil {
				return false, "", "", fmt.Errorf("failed to parse JSON Special Days value for: %v", contact)
			}
			for _, day := range days {
				var date time.Time
				if day.Date == "fathers-day-uk" {
					date, err = fathersday.FathersDay("uk", periodStart.Year())
					if err != nil {
						return false, "", "", fmt.Errorf("failed to get fathers day for: %v, %s", contact, err)
					}
				} else if day.Date == "mothering-sunday" {
					date, err = motheringsunday.MotheringSunday(periodStart.Year())
					if err != nil {
						return false, "", "", fmt.Errorf("failed to get mothering sunday for: %v, %s", contact, err)
					}
				} else {
					date, err = time.Parse("2006-01-02", day.Date)
					if err != nil {
						return false, "", "", fmt.Errorf("failed to parse special day date value for: %v, %s", contact, err)
					}
				}

				dateThisPeriod := time.Date(periodStart.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

				if dateInPeriod(periodStart, periodEnd, dateThisPeriod) {
					contactsWithSpecialDays = append(contactsWithSpecialDays, contact["Display Name"].(string))
					m := fmt.Sprintf("%s has a special day labelled '%s'", contact["Display Name"].(string), day.Label)
					if !today {
						m = fmt.Sprintf("%s on %s", m, date.Format("January 02"))
					}
					specialDayMessages = append(specialDayMessages, m)
				}
			}
		}
	}

	// handle when both are set
	if len(contactsWithSpecialDays) > 0 && len(contactsWithBirthdays) > 0 {
		allNames := append(contactsWithSpecialDays, contactsWithBirthdays...)
		allMessages := []string{}

		for _, message := range specialDayMessages {
			allMessages = append(allMessages, fmt.Sprintf("* %s", message))
		}
		for _, message := range birthdayMessages {
			allMessages = append(allMessages, fmt.Sprintf("* %s", message))
		}

		return true,
			fmt.Sprintf("%s have events", joinNamesList(allNames, ",", "&")),
			strings.Join(allMessages, "\n"),
			nil
	}

	// handle when only birthdays are set
	if len(contactsWithBirthdays) == 1 {
		return true,
			fmt.Sprintf("%s's Birthday", contactsWithBirthdays[0]),
			fmt.Sprintf("* %s", birthdayMessages[0]),
			nil
	}

	if len(contactsWithBirthdays) > 1 {
		return true,
			fmt.Sprintf("%d birthdays", len(contactsWithBirthdays)),
			fmt.Sprintf("* %s", strings.Join(birthdayMessages, "\n* ")),
			nil
	}

	// handle when only special days are set
	if len(specialDayMessages) == 1 {
		return true,
			specialDayMessages[0],
			specialDayMessages[0],
			nil
	}
	if len(specialDayMessages) > 1 {
		return true,
			fmt.Sprintf("%s have special days", joinNamesList(contactsWithSpecialDays, ",", "&")),
			fmt.Sprintf("* %s", strings.Join(specialDayMessages, "\n* ")),
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

func joinNamesList(names []string, separator, finalSeparator string) string {
	initialList := strings.Join(names[0:len(names)-1], separator+" ")

	return fmt.Sprintf("%s %s %s", initialList, finalSeparator, names[len(names)-1])
}
