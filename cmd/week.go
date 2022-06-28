package cmd

import (
	"fmt"
	"github.com/charlieegan3/airtable-contacts/pkg/webhook"
	"github.com/gomarkdown/markdown"
	"log"
	"time"

	"github.com/charlieegan3/airtable-contacts/pkg/airtable"
	"github.com/charlieegan3/airtable-contacts/pkg/pushover"
	"github.com/charlieegan3/airtable-contacts/pkg/specialdays"
	psh "github.com/gregdel/pushover"
	air "github.com/mehanizm/airtable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "send events for the current week",
	Run: func(cmd *cobra.Command, args []string) {
		// get the latest data
		airtableClient := air.NewClient(viper.GetString("airtable.key"))
		records, err := airtable.Download(
			airtableClient,
			viper.GetString("airtable.base"),
			viper.GetString("airtable.table"),
			viper.GetString("airtable.view"),
		)
		if err != nil {
			log.Fatalf("failed to download contacts: %s", err)
		}

		// set the notification period
		periodStart := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
		alert, title, body, err := specialdays.Generate(records, periodStart, 14, false)
		if err != nil {
			log.Fatalf("failed to generate alert body: %s", err)
		}

		// send the alert if needed
		pushoverRecipient := psh.NewRecipient(viper.GetString("pushover.user_key"))
		pushoverApp := psh.New(viper.GetString("pushover.app_token"))
		if alert {
			title := fmt.Sprintf("Weekly Summary (%s)", title)

			// notify on pushover
			err := pushover.Notify(
				pushoverApp,
				pushoverRecipient,
				title,
				body,
			)
			if err != nil {
				log.Fatalf("failed to send body: %s", err)
			}

			// notify via webhook
			err = webhook.Send(
				viper.GetString("webhook.endpoint"),
				title,
				string(markdown.ToHTML([]byte(body), nil, nil)),
				"https://airtable.com",
			)
			if err != nil {
				log.Fatalf("failed to send notification via webhook")
			}
		} else {
			title := "Weekly Summary (No Events)"
			body := "There are no events in the next two weeks"
			// notify on pushover
			err = pushover.Notify(
				pushoverApp,
				pushoverRecipient,
				title,
				body,
			)
			if err != nil {
				log.Fatalf("failed to send body: %s", err)
			}

			// notify via webhook
			err = webhook.Send(
				viper.GetString("webhook.endpoint"),
				title,

				string(markdown.ToHTML([]byte(body), nil, nil)),
				"https://airtable.com",
			)
			if err != nil {
				log.Fatalf("failed to send notification via webhook")
			}
		}
	},
}

func init() {
	notifyCmd.AddCommand(weekCmd)
}
