package carddav

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"

	govcard "github.com/emersion/go-vcard"
	"github.com/pkg/errors"

	"github.com/charlieegan3/tool-airtable-contacts/pkg/vcard"
)

func Sync(client Client, records []map[string]interface{}, vCardV3 bool, vCardPhotoSize int) error {

	reOrderPhotoAttrs := regexp.MustCompile(`(TYPE=[A-Z]+);(ENCODING=\w)`)
	rePhotoEncodingAttr := regexp.MustCompile(`ENCODING=(\w);`)
	rePhoneAttr := regexp.MustCompile(`TEL;TYPE=(\w+)`)

	existingItemsAll, err := client.GetAll()
	if err != nil {
		return errors.Wrap(err, "failed to list current vcards")
	}

	existingItemsCards, err := vcard.Parse(existingItemsAll)
	if err != nil {
		return errors.Wrap(err, "failed to parse current vcards")
	}

	recordsByID := make(map[string]map[string]interface{})
	for _, record := range records {
		id, ok := record["ID"].(string)
		if !ok {
			return errors.New("failed to get ID from record")
		}
		recordsByID[id] = record
	}

	doneIDs := make(map[string]bool)
	for _, existingCard := range existingItemsCards {
		id := existingCard.Get(govcard.FieldUID)
		if id == nil {
			return errors.New("failed to get ID from vcard")
		}
		idString := id.Value
		doneIDs[idString] = true

		expectedRecord, ok := recordsByID[idString]
		if !ok {
			err = client.Delete(idString)
			if err != nil {
				return errors.Wrap(err, "failed to delete outdated vcard")
			}
			continue
		}

		expectedVcardString, err := vcard.Generate(
			[]map[string]interface{}{expectedRecord},
			vCardV3,
			vCardPhotoSize,
			idString,
		)
		if err != nil {
			return errors.Wrap(err, "failed to generate new vcard")
		}

		buf := bytes.NewBufferString("")
		enc := govcard.NewEncoder(buf)
		err = enc.Encode(existingCard)
		if err != nil {
			return errors.Wrap(err, "failed to encode current vcard")
		}

		existingVcardString := strings.TrimSpace(buf.String())

		existingVcardString = reOrderPhotoAttrs.ReplaceAllString(existingVcardString, "$2;$1")
		existingVcardString = rePhotoEncodingAttr.ReplaceAllString(existingVcardString, "ENCODING=b;")
		existingVcardString = rePhoneAttr.ReplaceAllStringFunc(existingVcardString, func(s string) string {
			parts := rePhoneAttr.FindStringSubmatch(s)
			return fmt.Sprintf("TEL;TYPE=%s", strings.ToLower(parts[1]))
		})

		existingLines := strings.Split(strings.TrimSpace(existingVcardString), "\n")
		sort.Slice(existingLines, func(i, j int) bool {
			return existingLines[i] < existingLines[j]
		})
		expectedLines := strings.Split(strings.TrimSpace(expectedVcardString), "\n")
		sort.Slice(expectedLines, func(i, j int) bool {
			return expectedLines[i] < expectedLines[j]
		})

		shouldUpdate := false
		if len(existingLines) != len(expectedLines) {
			shouldUpdate = true
		}
		if strings.Join(existingLines, "") != strings.Join(expectedLines, "") {
			shouldUpdate = true
		}

		if shouldUpdate {
			fmt.Println("updating", idString)

			if len(existingLines) == len(expectedLines) {
				for i, line := range existingLines {
					if line != expectedLines[i] {
						fmt.Println("line", i, "differs")
						fmt.Println("old", firstN(line, 100))
						fmt.Println("new", firstN(expectedLines[i], 100))
					}
				}
			}

			err = client.Put(idString, expectedVcardString)
			if err != nil {
				return errors.Wrap(err, "failed to update vcard")
			}
		}
	}

	for id, rec := range recordsByID {
		if doneIDs[id] {
			continue
		}

		fmt.Println("creating", id)

		expectedVcardString, err := vcard.Generate(
			[]map[string]interface{}{rec},
			vCardV3,
			vCardPhotoSize,
			id,
		)

		err = client.Put(id, expectedVcardString)
		if err != nil {
			return errors.Wrap(err, "failed to update vcard")
		}
	}

	return nil
}

func firstN(str string, n int) string {
	v := []rune(str)
	if n >= len(v) {
		return str
	}
	return string(v[:n])
}
