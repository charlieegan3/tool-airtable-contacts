package cmd

import (
	"fmt"
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
		alert, title, message, err := specialdays.Generate(records, periodStart, 14)
		if err != nil {
			log.Fatalf("failed to generate alert message: %s", err)
		}

		// send the alert if needed
		pushoverRecipient := psh.NewRecipient(viper.GetString("pushover.user_key"))
		pushoverApp := psh.New(viper.GetString("pushover.app_token"))
		if alert {
			pushover.Notify(pushoverApp, pushoverRecipient, fmt.Sprintf("Weekly Summary (%s)", title), message)
		} else {
			pushover.Notify(pushoverApp, pushoverRecipient, "No Events", "There are no events in the next two weeks")
		}
	},
}

func init() {
	notifyCmd.AddCommand(weekCmd)
}
