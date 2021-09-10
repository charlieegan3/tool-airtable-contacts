/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"

	air "github.com/mehanizm/airtable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/charlieegan3/airtable-contacts/pkg/airtable"
	"github.com/charlieegan3/airtable-contacts/pkg/carddav"
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
		err = carddav.Sync(cardDAVClient, records)
		if err != nil {
			log.Fatalf("failed to upload to carddav: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
