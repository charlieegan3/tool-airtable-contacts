package vcard

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	govcard "github.com/emersion/go-vcard"
)

func Generate(contacts []map[string]interface{}, photoSize int) (string, error) {
	buf := bytes.NewBufferString("")
	enc := govcard.NewEncoder(buf)

	for _, contact := range contacts {
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
		for _, email := range computeEmailValues(contact) {
			card.AddValue(govcard.FieldEmail, email)
		}

		// set the phone numbers
		if val, ok := contact["Home Phone"].(string); ok {
			card.Add(govcard.FieldTelephone, &govcard.Field{
				Value: val,
				Params: map[string][]string{
					govcard.ParamType: {
						govcard.TypeHome,
					},
				},
			})
		}
		if val, ok := contact["Mobile Phone"].(string); ok {
			card.Add(govcard.FieldTelephone, &govcard.Field{
				Value: val,
				Params: map[string][]string{
					govcard.ParamType: {
						govcard.TypeCell,
					},
				},
			})
		}

		// set note
		if val, ok := contact["Note"].(string); ok {
			card.SetValue(govcard.FieldNote, val)
		}

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
						resp, err := http.Get(photoURL)
						if err != nil {
							return "", fmt.Errorf("failed to get photo: %s", err)
						}
						defer resp.Body.Close()

						image, _, err := image.Decode(resp.Body)
						if err != nil {
							return "", fmt.Errorf("failed to decode photo: %s", err)
						}
						dstImage := imaging.Fill(image, photoSize, photoSize, imaging.Center, imaging.Lanczos)

						buf := new(bytes.Buffer)
						err = jpeg.Encode(buf, dstImage, nil)
						if err != nil {
							return "", fmt.Errorf("failed to encode photo to jpg: %s", err)
						}

						card.Set(govcard.FieldPhoto, &govcard.Field{
							Value: base64.StdEncoding.EncodeToString(buf.Bytes()),
							Params: map[string][]string{
								govcard.ParamType: {"image/jpg"},
								"ENCODING":        {"b"},
							},
						})
					}
				}
			}
		}

		// write the card to output
		govcard.ToV4(card)
		err = enc.Encode(card)
		if err != nil {
			return "", fmt.Errorf("failed to encode vcard: %s", err)
		}
	}

	return strings.TrimSpace(buf.String()), nil
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

func computeEmailValues(contact map[string]interface{}) []string {
	emailString, ok := contact["Emails"].(string)

	if !ok {
		return []string{}
	}

	var emails []string
	for _, email := range strings.Split(emailString, "\n") {
		emails = append(emails, strings.TrimSpace(email))
	}

	return emails
}
