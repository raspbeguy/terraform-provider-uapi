package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func testClient(url string) *Client {
	return New(url, "tok", false)
}

func TestGetObjectFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer tok" {
			t.Errorf("missing bearer header, got %q", got)
		}
		w.Header().Set("ETag", `"abc123"`)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"r_1","managed":true,"enabled":true}`))
	}))
	defer srv.Close()

	obj, etag, found, err := testClient(srv.URL).GetObject(context.Background(), "/firewall/rules/r_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatal("expected found")
	}
	if etag != `"abc123"` {
		t.Errorf("etag = %q", etag)
	}
	if obj["id"] != "r_1" || obj["managed"] != true {
		t.Errorf("obj = %+v", obj)
	}
}

func TestGetObjectNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"not_found","message":"nope"}`))
	}))
	defer srv.Close()

	_, _, found, err := testClient(srv.URL).GetObject(context.Background(), "/x/y")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected not found")
	}
}

func TestValidationErrorEnvelope(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"code":"validation_failed","message":"bad","errors":[{"field":"target","code":"required","message":"is required"}]}`))
	}))
	defer srv.Close()

	_, _, err := testClient(srv.URL).Post(context.Background(), "/firewall/rules", map[string]any{}, "")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Code != "validation_failed" {
		t.Errorf("code = %q", apiErr.Code)
	}
	if len(apiErr.Errors) != 1 || apiErr.Errors[0].Field != "target" {
		t.Errorf("field errors not parsed: %+v", apiErr.Errors)
	}
}

func TestIfMatchSentAndETagReturned(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// uapi reads If-Match from the header OR the ?if_match= query param; assert both carry it.
		if r.URL.Query().Get("if_match") != `"v1"` {
			t.Errorf("if_match query = %q", r.URL.Query().Get("if_match"))
		}
		if r.Header.Get("If-Match") != `"v1"` {
			t.Errorf("If-Match header = %q", r.Header.Get("If-Match"))
		}
		w.Header().Set("ETag", `"v2"`)
		_, _ = w.Write([]byte(`{"id":"r_1"}`))
	}))
	defer srv.Close()

	_, etag, err := testClient(srv.URL).Put(context.Background(), "/firewall/rules/r_1", map[string]any{"target": "ACCEPT"}, `"v1"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if etag != `"v2"` {
		t.Errorf("new etag = %q", etag)
	}
}

func TestPreconditionFailed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPreconditionFailed)
		_, _ = w.Write([]byte(`{"code":"precondition_failed","message":"stale"}`))
	}))
	defer srv.Close()

	_, _, err := testClient(srv.URL).Put(context.Background(), "/firewall/rules/r_1", map[string]any{}, `"stale"`)
	if !IsPreconditionFailed(err) {
		t.Fatalf("expected precondition-failed, got %v", err)
	}
}

func TestRetryOnLocked(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n < 3 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusLocked)
			_, _ = w.Write([]byte(`{"code":"locked","message":"busy"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"r_1"}`))
	}))
	defer srv.Close()

	obj, _, err := testClient(srv.URL).Put(context.Background(), "/firewall/rules/r_1", map[string]any{"target": "ACCEPT"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj["id"] != "r_1" {
		t.Errorf("id = %v", obj["id"])
	}
	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Errorf("expected 3 calls (2 locked + 1 ok), got %d", got)
	}
}

func TestLockedExhaustsRetries(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusLocked)
		_, _ = w.Write([]byte(`{"code":"locked","message":"busy"}`))
	}))
	defer srv.Close()

	_, _, err := testClient(srv.URL).Post(context.Background(), "/firewall/rules", map[string]any{}, "")
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	if apiErr, ok := err.(*APIError); !ok || apiErr.Status != http.StatusLocked {
		t.Errorf("expected locked APIError, got %v", err)
	}
}

func TestDeleteToleratesNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"not_found"}`))
	}))
	defer srv.Close()

	if err := testClient(srv.URL).Delete(context.Background(), "/firewall/rules/gone", ""); err != nil {
		t.Fatalf("delete should tolerate 404, got %v", err)
	}
}

func TestGetList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"mac":"aa:bb:cc:dd:ee:ff","ip":"10.0.0.2"}]`))
	}))
	defer srv.Close()

	list, err := testClient(srv.URL).GetList(context.Background(), "/dhcp/leases")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 || list[0]["ip"] != "10.0.0.2" {
		t.Errorf("unexpected list: %+v", list)
	}
}
