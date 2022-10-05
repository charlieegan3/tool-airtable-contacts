package cmd

import (
	"log"

	air "github.com/mehanizm/airtable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/charlieegan3/tool-airtable-contacts/pkg/airtable"
	"github.com/charlieegan3/tool-airtable-contacts/pkg/carddav"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "download data from airtable and upload to carddav server",
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

		cardDAVClient := carddav.Client{
			URL:      viper.GetString("carddav.serverURL"),
			User:     viper.GetString("carddav.user"),
			Password: viper.GetString("carddav.password"),
		}

		// records passed as vcard sync is done on per contact basis
		err = carddav.Sync(cardDAVClient, records, viper.GetBool("vcard.use_v3"), viper.GetInt("vcard.photos.size"))
		if err != nil {
			log.Fatalf("failed to upload to carddav: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
