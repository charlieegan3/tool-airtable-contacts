package vcard

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"image/jpeg"
	"net/http"
	"regexp"
	"strings"

	"github.com/disintegration/imaging"
	govcard "github.com/emersion/go-vcard"
)

func CRC32Hash(input any) string {
	str := fmt.Sprintf("%v", input)
	return fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(str)))
}

func Generate(contact map[string]interface{}, useV3 bool, photoSize int, id string) (string, error) {

	buf := bytes.NewBufferString("")
	enc := govcard.NewEncoder(buf)

	card := make(govcard.Card)

	// set the name and formatted name
	displayName, firstName, lastName, err := computeNameValues(contact)
	if err != nil {
		return "", fmt.Errorf("failed to compute names: %s", err)
	}

	card.SetValue(govcard.FieldFormattedName, displayName)
	card.AddName(&govcard.Name{
		GivenName:  firstName,
		FamilyName: lastName,
	})

	// set the emails
	if val, ok := contact["JSON Emails"].(string); ok {
		decoder := json.NewDecoder(strings.NewReader(val))
		var emails []struct {
			Label     string `json:"label"`
			Value     string `json:"value"`
			Preferred bool   `json:"preferred"`
		}
		err := decoder.Decode(&emails)
		if err != nil {
			return "", fmt.Errorf("failed to parse JSON emails: %s", err)
		}
		for i, email := range emails {
			preferred := "1"
			if email.Preferred {
				preferred = "2"
			}
			card.Add(fmt.Sprintf("item%d.%s", i+1, govcard.FieldEmail), &govcard.Field{
				Value: email.Value,
				Params: map[string][]string{
					govcard.ParamPreferred: {preferred},
				},
			})
			card.Add(fmt.Sprintf("item%d.X-ABLABEL", i+1), &govcard.Field{
				Value: email.Label,
			})
		}
	}

	// set the phone numbers
	if val, ok := contact["JSON Phone Numbers"].(string); ok {
		decoder := json.NewDecoder(strings.NewReader(val))
		var numbers []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}
		err := decoder.Decode(&numbers)
		if err != nil {
			return "", fmt.Errorf("failed to parse JSON phone numbers: %s", err)
		}
		for _, number := range numbers {
			card.Add(govcard.FieldTelephone, &govcard.Field{
				Value: number.Value,
				Params: map[string][]string{
					govcard.ParamType: {number.Type},
				},
			})
		}
	}

	// set addresses
	if val, ok := contact["JSON Addresses"].(string); ok {
		decoder := json.NewDecoder(strings.NewReader(val))
		var addresses []govcard.Address
		err := decoder.Decode(&addresses)
		if err != nil {
			return "", fmt.Errorf("failed to parse JSON addresses: %s", err)
		}
		for _, address := range addresses {
			card.AddAddress(&address)
		}
	}

	// set note
	var noteText string
	if val, ok := contact["Note"].(string); ok {
		noteText = val
	}
	noteText = fmt.Sprintf("%s\n%s", noteText, CRC32Hash(contact))
	card.SetValue(govcard.FieldNote, noteText)

	// set org from company value
	if val, ok := contact["Company"].(string); ok {
		card.SetValue(govcard.FieldOrganization, val)
	}

	// set photo value
	if val, ok := contact["Profile Image"].([]interface{}); ok {
		if len(val) > 0 {
			photo, ok := val[0].(map[string]interface{})
			if ok {
				photoURL, photoURLOk := photo["url"].(string)
				if photoURLOk {
					buf, err := fetchPhoto(photoURL, photoSize)
					if err != nil {
						return "", err
					}

					card.Set(govcard.FieldPhoto, &govcard.Field{
						Value: base64.StdEncoding.EncodeToString(buf.Bytes()),
						Params: map[string][]string{
							govcard.ParamType: {"JPEG"},
							"ENCODING":        {"b"},
						},
					})
				}
			}
		}
	}

	// optionally, set the ID
	if id != "" {
		card.SetValue("UID", id)
	}

	// allow setting of the version
	if useV3 {
		card.SetValue(govcard.FieldVersion, "3.0")
	} else {
		govcard.ToV4(card)
	}

	// write the card to output
	err = enc.Encode(card)
	if err != nil {
		return "", fmt.Errorf("failed to encode vcard: %s", err)
	}

	rePhotoEncodingAttr := regexp.MustCompile(`ENCODING=\w;`)
	return rePhotoEncodingAttr.ReplaceAllString(buf.String(), "ENCODING=B;"), nil
}

func fetchPhoto(photoURL string, photoSize int) (*bytes.Buffer, error) {
	resp, err := http.Get(photoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get photo: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get photo: %s", resp.Status)
	}

	img, err := imaging.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode photo: %s", err)
	}
	dstImage := imaging.Fill(img, photoSize, photoSize, imaging.Center, imaging.Lanczos)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, dstImage, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode photo to jpg: %s", err)
	}
	return buf, nil
}

func computeNameValues(contact map[string]interface{}) (string, string, string, error) {
	displayName, ok := contact["Display Name"].(string)

	if !ok || displayName == "" {
		return "", "", "", fmt.Errorf("contact was missing Display Name")
	}

	words := strings.Split(strings.TrimSpace(displayName), " ")

	switch len(words) {
	case 1:
		return displayName, displayName, "", nil
	case 2:
		return displayName, words[0], words[1], nil
	default:
		return displayName, words[0], strings.Join(words[1:], " "), nil
	}
}
