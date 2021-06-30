package vcard

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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
		if val, ok := contact["JSON Emails"].(string); ok {
			decoder := json.NewDecoder(strings.NewReader(val))
			emails := []struct {
				Label     string `json:"label"`
				Value     string `json:"value"`
				Preferred bool   `json:"preferred"`
			}{}
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
			numbers := []struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			}{}
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
