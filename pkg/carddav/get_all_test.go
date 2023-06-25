package carddav

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestGetAll(t *testing.T) {
	testServerCalled := false
	body := `BEGIN:VCARD
VERSION:3.0
FN:Example
N:;Example;;;
NOTE:notes\n
UID:example_123
item1.EMAIL;PREF=2:example@example.com
item1.X-ABLABEL:Email
END:VCARD
BEGIN:VCARD
VERSION:3.0
FN:Example2
N:;Example;;;
NOTE:notes2\n
UID:example_1234
item1.EMAIL;PREF=2:example2@example.com
item1.X-ABLABEL:Email
END:VCARD`

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		td.Cmp(t, r.Method, "GET")
		td.Cmp(t, r.Header["Authorization"], []string{"Basic YWxpY2U6cGFzc3dvcmQ="})
		td.Cmp(t, r.URL.Path, "/dav/addressbooks/user/charlieegan3@example.com/Default")
		testServerCalled = true

		w.Write([]byte(body))
	}))

	cardDavClient := Client{
		URL:      testServer.URL + "/dav/addressbooks/user/charlieegan3@example.com/Default/",
		User:     "alice",
		Password: "password",
	}

	result, err := cardDavClient.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	td.Cmp(t, testServerCalled, true)

	td.Cmp(t, result, body)
}
