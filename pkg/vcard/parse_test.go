package vcard

import "testing"

func TestParse(t *testing.T) {
	input := `BEGIN:VCARD
VERSION:4.0
FN:John Appleseed
N:Appleseed;John;;;
END:VCARD
BEGIN:VCARD
VERSION:4.0
FN:Jane Appleseed
N:Appleseed;Jane;;;
END:VCARD`

	results, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if results[0].Name().GivenName != "John" {
		t.Fatalf("expected John, got %s", results[0].Name().GivenName)
	}
	if results[1].Name().GivenName != "Jane" {
		t.Fatalf("expected Jane, got %s", results[1].Name().GivenName)
	}
}
