package carddav

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestDelete(t *testing.T) {
	testServerCalled := false
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		td.Cmp(t, r.Method, "DELETE")
		td.Cmp(t, r.Header["Authorization"], []string{"Basic YWxpY2U6cGFzc3dvcmQ="})
		td.Cmp(t, r.URL.Path, "/dav/addressbooks/user/charlieegan3@fastmail.com/Default/test-id.vcf")
		testServerCalled = true
	}))

	cardDavClient := Client{
		URL:      testServer.URL + "/dav/addressbooks/user/charlieegan3@fastmail.com/Default/",
		User:     "alice",
		Password: "password",
	}

	err := cardDavClient.Delete("test-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	td.Cmp(t, testServerCalled, true)
}
