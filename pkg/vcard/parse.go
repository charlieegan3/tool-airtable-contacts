package vcard

import (
	"io"
	"strings"

	"github.com/emersion/go-vcard"
)

func Parse(input string) ([]vcard.Card, error) {
	var cards []vcard.Card

	dec := vcard.NewDecoder(strings.NewReader(input))
	for {
		card, err := dec.Decode()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}
