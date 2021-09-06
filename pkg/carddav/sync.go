package carddav

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/charlieegan3/airtable-contacts/pkg/vcard"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func Sync(client Client, records []map[string]interface{}) error {
	// get the existing list of cards
	existingItems, err := client.List()
	if err != nil {
		return errors.Wrap(err, "failed to list current vcards")
	}

	// first create any missing contacts, or new ones for updates
	currentItems := []string{}
	for _, record := range records {
		// hash to use as ID
		h := sha1.New()
		h.Write([]byte(fmt.Sprintf("%v", record)))
		id := hex.EncodeToString(h.Sum(nil))
		currentItems = append(currentItems, id)

		shouldSkip := false
		for _, v := range existingItems {
			if id == v {
				shouldSkip = true
			}
		}
		if shouldSkip {
			continue
		}

		vcardString, err := vcard.Generate(
			[]map[string]interface{}{record},
			viper.GetBool("vcard.use_v3"),
			viper.GetInt("vcard.photo.size"),
			id,
		)
		if err != nil {
			return errors.Wrap(err, "failed to generate new vcard")
		}

		err = client.Put(id, vcardString)
		if err != nil {
			return errors.Wrap(err, "failed to create new vcard")
		}
	}

	// finally, remove any old items not in the new list
	for _, existingItem := range existingItems {
		shouldDelete := true
		for _, currentItem := range currentItems {
			if currentItem == existingItem {
				shouldDelete = false
				break
			}
		}
		if shouldDelete {
			err = client.Delete(existingItem)
			if err != nil {
				return errors.Wrap(err, "failed to delete outdated vcard")
			}
		}
	}

	return nil
}
