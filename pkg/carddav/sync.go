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

		var ok bool
		expectedRecord, ok := recordsByID[idString]
		if !ok {
			err = client.Delete(idString)
			if err != nil {
				return errors.Wrap(err, "failed to delete outdated vcard")
			}
			continue
		}

		needsUpdate := false
		note := existingCard.Get(govcard.FieldNote)
		if note == nil {
			needsUpdate = true
		} else {
			lines := strings.Split(note.Value, "\n")
			var lastLine string
			if len(lines) > 0 {
				lastLine = strings.TrimSpace(lines[len(lines)-1])
			}
			if lastLine != vcard.CRC32Hash(expectedRecord) {
				needsUpdate = true
			}
		}

		if !needsUpdate {
			continue
		}

		err := updateVcard(
			client,
			expectedRecord,
			existingCard,
			vCardV3,
			vCardPhotoSize,
			idString,
		)
		if err != nil {
			return errors.Wrap(err, "failed to update vcard")
		}
	}

	for id, rec := range recordsByID {
		if doneIDs[id] {
			continue
		}

		fmt.Println("creating", id)

		expectedVcardString, err := vcard.Generate(
			rec,
			vCardV3,
			vCardPhotoSize,
			id,
		)
		expectedVcardString = normalizeVcard(expectedVcardString)

		err = client.Put(id, expectedVcardString)
		if err != nil {
			return errors.Wrap(err, "failed to update vcard")
		}
	}

	return nil
}

func updateVcard(
	client Client,
	expectedRecord map[string]interface{},
	existingCard govcard.Card,
	vCardV3 bool,
	vCardPhotoSize int,
	idString string,
) error {
	expectedVcardString, err := vcard.Generate(
		expectedRecord,
		vCardV3,
		vCardPhotoSize,
		idString,
	)
	if err != nil {
		return errors.Wrap(err, "failed to generate new vcard")
	}

	expectedVcardString = normalizeVcard(expectedVcardString)

	var existingVcardString string
	buf := bytes.NewBufferString("")
	enc := govcard.NewEncoder(buf)
	err = enc.Encode(existingCard)
	if err != nil {
		return errors.Wrap(err, "failed to encode current vcard")
	}

	existingVcardString = normalizeVcard(buf.String())

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

	return nil
}

func normalizeVcard(vcardString string) string {
	reOrderPhotoAttrs := regexp.MustCompile(`(ENCODING=\w);(TYPE=[A-Z]+)`)
	rePhoneAttr := regexp.MustCompile(`TEL;TYPE=(\w+)`)

	newVcardString := strings.TrimSpace(vcardString)

	newVcardString = reOrderPhotoAttrs.ReplaceAllString(newVcardString, "$2;$1")
	newVcardString = rePhoneAttr.ReplaceAllStringFunc(newVcardString, func(s string) string {
		parts := rePhoneAttr.FindStringSubmatch(s)
		return fmt.Sprintf("TEL;TYPE=%s", strings.ToLower(parts[1]))
	})

	newVcardString = strings.ReplaceAll(newVcardString, "ENCODING=b", "ENCODING=B")

	return newVcardString
}

func firstN(str string, n int) string {
	v := []rune(str)
	if n >= len(v) {
		return str
	}
	return string(v[:n])
}
