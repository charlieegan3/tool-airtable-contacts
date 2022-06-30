package cmd

import (
	"github.com/charlieegan3/airtable-contacts/pkg/webhook"
	"github.com/gomarkdown/markdown"
	"log"
	"time"

	psh "github.com/gregdel/pushover"
	air "github.com/mehanizm/airtable"
	"github.com/spf13/cobra"

	"github.com/charlieegan3/airtable-contacts/pkg/airtable"
	"github.com/charlieegan3/airtable-contacts/pkg/pushover"
	"github.com/charlieegan3/airtable-contacts/pkg/specialdays"
	"github.com/spf13/viper"
)

var dayCmd = &cobra.Command{
	Use:   "day",
	Short: "send events for current day",
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
		alert, title, body, err := specialdays.Generate(records, periodStart, 1, true)
		if err != nil {
			log.Fatalf("failed to generate alert message: %s", err)
		}

		// send the alert if needed
		if alert {
			// notify via pushover
			pushoverRecipient := psh.NewRecipient(viper.GetString("pushover.user_key"))
			pushoverApp := psh.New(viper.GetString("pushover.app_token"))
			err := pushover.Notify(pushoverApp, pushoverRecipient, title, body)
			if err != nil {
				log.Fatalf("failed to send notification via pushover")
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
	notifyCmd.AddCommand(dayCmd)
}
