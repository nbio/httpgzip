package httpgzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGzipResponse(t *testing.T) {
	server := createServer()
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Accept-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Errorf("expected gzip response; got headers: %v", resp.Header)
	}

	if resp.Header.Get("Vary") != "Accept-Encoding" {
		t.Errorf("expected Vary: Accept-Encoding; got headers: %v", resp.Header)
	}

	var message bytes.Buffer
	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(&message, gr)
	if message.String() != "hello" {
		t.Fatal("expected 'hello' to roundtrip, got: %q", message.String())
	}
}

func TestNonGzipResponse(t *testing.T) {
	server := createServer()
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") == "gzip" {
		t.Errorf("unexpected gzip response; headers: %v", resp.Header)
	}

	if resp.Header.Get("Vary") != "Accept-Encoding" {
		t.Errorf("expected Vary: Accept-Encoding; got headers: %v", resp.Header)
	}

	var message bytes.Buffer
	io.Copy(&message, resp.Body)
	if message.String() != "hello" {
		t.Fatal("expected 'hello' to roundtrip, got: %q", message.String())
	}
}

func createServer() *httptest.Server {
	h := GzipResponse(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("hello"))
	}))
	mux := http.NewServeMux()
	mux.Handle("/", h)
	return httptest.NewServer(mux)
}