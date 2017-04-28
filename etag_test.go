package etag

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acoshift/header"
)

func TestETag(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	ts := httptest.NewServer(New(DefaultConfig)(h))
	defer ts.Close()

	resp, _ := http.Get(ts.URL)
	et := resp.Header.Get(header.ETag)
	if len(et) == 0 {
		t.Fatalf("expected et is not empty string")
	}

	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Add(header.IfNoneMatch, et)
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusNotModified {
		t.Fatalf("expected response status to be 304; got %d", resp.StatusCode)
	}

	req, _ = http.NewRequest(http.MethodHead, ts.URL, nil)
	req.Header.Add(header.IfNoneMatch, et)
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusNotModified {
		t.Fatalf("expected response status to be 304; got %d", resp.StatusCode)
	}

	req, _ = http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Add(header.IfNoneMatch, "\"other_etag\", "+et)
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusNotModified {
		t.Fatalf("expected response status to be 304; got %d", resp.StatusCode)
	}

	req, _ = http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Add(header.IfNoneMatch, "*")
	resp, _ = http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusNotModified {
		t.Fatalf("expected response status to be 304; got %d", resp.StatusCode)
	}
}
