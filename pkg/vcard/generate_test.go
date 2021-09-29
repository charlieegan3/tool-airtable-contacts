package vcard

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
)

func TestGenerate(t *testing.T) {
	jpgImage, err := os.Open("example.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer jpgImage.Close()
	jpgImageData, err := ioutil.ReadAll(jpgImage)
	if err != nil {
		t.Fatal(err)
	}

	pngImage, err := os.Open("example.png")
	if err != nil {
		t.Fatal(err)
	}
	defer jpgImage.Close()
	pngImageData, err := ioutil.ReadAll(pngImage)
	if err != nil {
		t.Fatal(err)
	}

	profilePhotoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/jpg":
			w.Write(jpgImageData)
		case "/png":
			w.Write(pngImageData)
		default:
			t.Fatalf("unexpected path called in test image server")
		}
	}))

	testCases := []struct {
		Description    string
		Fields         []map[string]interface{}
		ExpectedOutput string
		// used when the field label order is not deterministic
		ExpectedContains string
		UseV3            bool
		ID               string
	}{
		{
			Description: "it sets the name correctly",
			Fields: []map[string]interface{}{
				{"Display Name": "John Appleseed"},
			},
			ExpectedOutput: `BEGIN:VCARD
VERSION:4.0
FN:John Appleseed
N:Appleseed;John;;;
END:VCARD`,
		},
		{
			Description: "it can generate a v3 card",
			UseV3:       true,
			Fields: []map[string]interface{}{
				{"Display Name": "John Appleseed"},
			},
			ExpectedOutput: `BEGIN:VCARD
VERSION:3.0
FN:John Appleseed
N:Appleseed;John;;;
END:VCARD`,
		},
		{
			Description: "it can create multiple cards",
			Fields: []map[string]interface{}{
				{"Display Name": "John Appleseed"},
				{"Display Name": "Jane Appleseed"},
			},
			ExpectedOutput: `BEGIN:VCARD
VERSION:4.0
FN:John Appleseed
N:Appleseed;John;;;
END:VCARD
BEGIN:VCARD
VERSION:4.0
FN:Jane Appleseed
N:Appleseed;Jane;;;
END:VCARD`,
		},
		{
			Description: "it sets emails",
			Fields: []map[string]interface{}{
				{
					"Display Name": "John Appleseed",
					"JSON Emails":  `[{"label":"Personal", "value":"xxxxxxx@xxxxxxx.uk", "preferred": true},{"label":"Email", "value":"xxxx@xxxxx.info", "preferred": false}]`,
				},
			},
			ExpectedOutput: `BEGIN:VCARD
VERSION:4.0
FN:John Appleseed
N:Appleseed;John;;;
item1.EMAIL;PREF=2:xxxxxxx@xxxxxxx.uk
item1.X-ABLABEL:Personal
item2.EMAIL;PREF=1:xxxx@xxxxx.info
item2.X-ABLABEL:Email
END:VCARD`,
		},
		{
			Description: "it sets phone numbers",
			Fields: []map[string]interface{}{
				{
					"Display Name":       "John Appleseed",
					"JSON Phone Numbers": `[{"type":"mobile", "value":"+44333 666 7777"},{"type":"home", "value":"01526 555555"}]`,
				},
			},
			ExpectedOutput: `BEGIN:VCARD
VERSION:4.0
FN:John Appleseed
N:Appleseed;John;;;
TEL;TYPE=mobile:+44333 666 7777
TEL;TYPE=home:01526 555555
END:VCARD`,
		},
		{
			Description: "it sets notes",
			Fields: []map[string]interface{}{
				{
					"Display Name": "John Appleseed",
					"Note":         "test",
				},
			},
			ExpectedOutput: `BEGIN:VCARD
VERSION:4.0
FN:John Appleseed
N:Appleseed;John;;;
NOTE:test
END:VCARD`,
		},
		{
			Description: "it sets org values for companies",
			Fields: []map[string]interface{}{
				{
					"Display Name": "John Appleseed",
					"Company":      "Apple",
				},
			},
			ExpectedOutput: `BEGIN:VCARD
VERSION:4.0
FN:John Appleseed
N:Appleseed;John;;;
ORG:Apple
END:VCARD`,
		},
		{
			Description: "it sets png images",
			Fields: []map[string]interface{}{
				{
					"Display Name": "John Appleseed",
					"Profile Image": []interface{}{
						map[string]interface{}{
							"type": "image/png",
							"url":  profilePhotoServer.URL + "/png",
						},
					},
				},
			},
			ExpectedContains: `/9j/2wCEAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDIBCQkJDAsMGA0NGDIhHCEyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMv/AABEIABkAGQMBIgACEQEDEQH/xAGiAAABBQEBAQEBAQAAAAAAAAAAAQIDBAUGBwgJCgsQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29/j5+gEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoLEQACAQIEBAMEBwUEBAABAncAAQIDEQQFITEGEkFRB2FxEyIygQgUQpGhscEJIzNS8BVictEKFiQ04SXxFxgZGiYnKCkqNTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqCg4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2dri4+Tl5ufo6ery8/T19vf4+fr/2gAMAwEAAhEDEQA/AOA07w1bW2jWtwy21wLpkxMxVlG4dDkZXr1A7ViSeG45vFP2SETxac8nMu3f5SnJwSBjt149aZoPivVNOjg01LrbZGZWAKgmM7gcq3VeeeK9NhjnaRomtZ724jctGiEIoBHO7HJ+bjHGRxnpW+KxUHSjTjGz7/1udOGwntk25WOZ1Lwxo0mj3MUEdvb3duCI2DuHLKPm37uxIPbgkc15rXWeMdTvItSudKWWBbZWDFLfGORnG4dcZrk64qKkl7zOONGdKUoTd7MK9PsfHWhpp1tLcpcm6WFY5YNm5H28+oByfXOK8woqqlNVFqbQqSh8LsSTy+fcSS4xvYtjPqajoorQg//Z`,
		},
		{
			Description: "it sets jpg images",
			Fields: []map[string]interface{}{
				{
					"Display Name": "John Appleseed",
					"Profile Image": []interface{}{
						map[string]interface{}{
							"type": "image/jpg",
							"url":  profilePhotoServer.URL + "/jpg",
						},
					},
				},
			},
			ExpectedContains: `/9j/2wCEAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDIBCQkJDAsMGA0NGDIhHCEyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMv/AABEIABkAGQMBIgACEQEDEQH/xAGiAAABBQEBAQEBAQAAAAAAAAAAAQIDBAUGBwgJCgsQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29/j5+gEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoLEQACAQIEBAMEBwUEBAABAncAAQIDEQQFITEGEkFRB2FxEyIygQgUQpGhscEJIzNS8BVictEKFiQ04SXxFxgZGiYnKCkqNTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqCg4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2dri4+Tl5ufo6ery8/T19vf4+fr/2gAMAwEAAhEDEQA/APHo4hKVRQWbsAMk13HgzwvcW2qLq98rQR2y7o4mHzSFgQDjsuD/ACpnw902HzptQnxuQhIycjaepI9+g/Ouy1C8+xRi6gVSq5WUbuQD0z+OKybd7IqMdLs5PxpNp8ujm0tHjfFwrFo02jdggj361wX2b3ro9e1JbmZrfAXZJufjo2MEfzrEyPWmtCXqzrtIa8TQpJYrc+T52Fk6Ddjke9aMN3eXGnXaQ6beyssYUKsDYkJPOCRjgfnXUeEP+QRafRf610i/6kfVqydTlexoo6Hmh8Ate3Ml3cTyQRTuZFGxUZcn7rBjwfpxTv8AhXVj/wBBKX/xyurn+5c/7w/9BFc1Ue0Y+VI//9k=`,
		},
		{
			Description: "it sets the id if present",
			Fields: []map[string]interface{}{
				{"Display Name": "John Appleseed"},
			},
			ID: "example",
			ExpectedOutput: `BEGIN:VCARD
VERSION:4.0
FN:John Appleseed
N:Appleseed;John;;;
UID:example
END:VCARD`,
		},
		{
			Description: "it sets addresses correctly",
			Fields: []map[string]interface{}{
				{
					"Display Name": "John Appleseed",
					"JSON Addresses": `[
				      {"ExtendedAddress":"","StreetAddress":"15 Whitby Road","Region":"London","PostalCode":"IV65NE","Country":"UK"},
				      {"ExtendedAddress":"Flat 2","StreetAddress":"Sweden Road","Region":"Vaxjo","PostalCode":"12345","Country":"Sweden"}]`,
				},
			},
			ExpectedOutput: `BEGIN:VCARD
VERSION:4.0
ADR:;;15 Whitby Road;;London;IV65NE;UK
ADR:;Flat 2;Sweden Road;;Vaxjo;12345;Sweden
FN:John Appleseed
N:Appleseed;John;;;
END:VCARD`,
		},
	}

	for _, test := range testCases {
		t.Run(test.Description, func(t *testing.T) {
			// 25 is used here so the base64 data is small in the examples above
			result, err := Generate(test.Fields, test.UseV3, 25, test.ID)
			// re-format generated string with \r for comparison in tests
			result = strings.ReplaceAll(result, "\r", "")
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if test.ExpectedOutput != "" && result != test.ExpectedOutput {
				t.Fatalf("unexpected result:\n%s", diff.LineDiff(result, test.ExpectedOutput))
			}
			if test.ExpectedContains != "" && !strings.Contains(result, test.ExpectedContains) {
				t.Fatalf("expected substring not found:\n%s in result: %s", test.ExpectedContains, result)
			}
		})
	}
}
