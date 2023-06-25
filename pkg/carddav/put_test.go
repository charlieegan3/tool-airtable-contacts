package carddav

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestPut(t *testing.T) {
	expectedVCard := `BEGIN:VCARD
VERSION:3.0
UID:charlie
N:;Charlie;;;
FN:Charlie
NOTE:
END:VCARD`

	testServerCalled := false
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		td.Cmp(t, r.Method, "PUT")
		td.Cmp(t, r.Header["Authorization"], []string{"Basic YWxpY2U6cGFzc3dvcmQ="})
		td.Cmp(t, r.URL.Path, "/dav/addressbooks/user/charlieegan3@example.com/Default/test-id.vcf")

		bytes, err := ioutil.ReadAll(r.Body)
		td.Cmp(t, err, nil)
		td.Cmp(t, string(bytes), expectedVCard)

		testServerCalled = true
	}))

	cardDavClient := Client{
		URL:      testServer.URL + "/dav/addressbooks/user/charlieegan3@example.com/Default/",
		User:     "alice",
		Password: "password",
	}

	err := cardDavClient.Put("test-id", expectedVCard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	td.Cmp(t, testServerCalled, true)
}
