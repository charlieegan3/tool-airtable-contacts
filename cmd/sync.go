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
		viper.SetConfigName("config.test")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Fatal error config file: %s \n", err)
		}

		// get the latest data
		airtableClient := air.NewClient(
			viper.GetString("tools.airtable-contacts.airtable.key"),
		)
		records, err := airtable.Download(
			airtableClient,
			viper.GetString("tools.airtable-contacts.airtable.base"),
			viper.GetString("tools.airtable-contacts.airtable.table"),
			viper.GetString("tools.airtable-contacts.airtable.view"),
		)
		if err != nil {
			log.Fatalf("failed to download contacts: %s", err)
		}

		cardDAVClient := carddav.Client{
			URL:      viper.GetString("tools.airtable-contacts.carddav.server_url"),
			User:     viper.GetString("tools.airtable-contacts.carddav.user"),
			Password: viper.GetString("tools.airtable-contacts.carddav.password"),
		}

		// records passed as vcard sync is done on per contact basis
		err = carddav.Sync(
			cardDAVClient,
			records,
			viper.GetBool("tools.airtable-contacts.vcard.use_v3"),
			viper.GetInt("tools.airtable-contacts.vcard.photo.size"),
		)
		if err != nil {
			log.Fatalf("failed to upload to carddav: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
