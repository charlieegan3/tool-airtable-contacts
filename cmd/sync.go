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
	"os"

	dbx "github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	air "github.com/mehanizm/airtable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/charlieegan3/airtable-contacts/pkg/airtable"
	"github.com/charlieegan3/airtable-contacts/pkg/carddav"
	"github.com/charlieegan3/airtable-contacts/pkg/dropbox"
	"github.com/charlieegan3/airtable-contacts/pkg/vcard"
)

var syncDropbox, syncCardDAV, syncFile bool

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "download data from airtable and upload to dropbox",
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

		// generate string for all contacts
		vcardString, err := vcard.Generate(
			records,
			viper.GetBool("vcard.use_v3"),
			viper.GetInt("vcard.photo.size"),
			"",
		)
		if err != nil {
			log.Fatal(err)
		}

		if syncCardDAV {
			cardDAVClient := carddav.Client{
				URL:      viper.GetString("carddav.serverURL"),
				User:     viper.GetString("carddav.user"),
				Password: viper.GetString("carddav.password"),
			}

			// records passed as vcard sync is done on per contact basis
			err := carddav.Sync(cardDAVClient, records)
			if err != nil {
				log.Fatal(err)
			}
		}
		if syncDropbox {
			dropboxClient := files.New(dbx.Config{
				Token:    viper.GetString("dropbox.token"),
				LogLevel: dbx.LogOff,
			})
			dropbox.Upload(dropboxClient, viper.GetString("dropbox.path"), []byte(vcardString))
		}
		if syncFile {
			err = os.WriteFile("out.vcard", []byte(vcardString), 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	syncCmd.Flags().BoolVar(&syncDropbox, "dropbox", false, "if set, dropbox will be synced")
	syncCmd.Flags().BoolVar(&syncCardDAV, "carddav", false, "if set, carddav will be synced")
	syncCmd.Flags().BoolVar(&syncFile, "file", false, "if set, local will saved")
	rootCmd.AddCommand(syncCmd)
}
