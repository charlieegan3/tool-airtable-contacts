package carddav

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestBackoff(t *testing.T) {
	requestShouldFail := true
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestShouldFail {
			requestShouldFail = false
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		td.Cmp(t, r.Method, "PROPFIND")
		td.Cmp(t, r.Header["Authorization"], []string{"Basic YWxpY2U6cGFzc3dvcmQ="})
		td.Cmp(t, r.Header["Depth"], []string{"1"})

		b, err := ioutil.ReadFile("propfind.xml")
		td.Cmp(t, err, nil)
		w.Write(b)
	}))

	expectedItems := []string{
		"6316e64d-176c-42bc-9d59-30ed53c1a06b",
		"a5fadfc4-27fc-459c-847d-c8dabe3f789e",
	}

	cardDavClient := Client{
		URL:      testServer.URL + "/dav/addressbooks/user/charlieegan3@fastmail.com/Default",
		User:     "alice",
		Password: "password",
	}

	gotItems, err := cardDavClient.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	td.Cmp(t, gotItems, expectedItems)
}
