package render

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAudioContentType(t *testing.T) {
	tests := []struct {
		filename string
		want     string
	}{
		{"narration.mp3", "audio/mpeg"},
		{"narration.WAV", "audio/wav"},
		{"a.ogg", "audio/ogg"},
		{"a.m4a", "audio/aac"},
		{"a.aac", "audio/aac"},
		{"a.flac", "audio/flac"},
		{"unknown.bin", "application/octet-stream"},
		{"noext", "application/octet-stream"},
	}
	for _, tt := range tests {
		if got := AudioContentType(tt.filename); got != tt.want {
			t.Errorf("AudioContentType(%q) = %q, want %q", tt.filename, got, tt.want)
		}
	}
}

func TestDownloadURL(t *testing.T) {
	content := []byte("fake-mp4-bytes")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(content)
	}))
	defer srv.Close()

	var buf bytes.Buffer
	if err := DownloadURL(context.Background(), nil, srv.URL, &buf); err != nil {
		t.Fatalf("DownloadURL() error = %v", err)
	}
	if !bytes.Equal(buf.Bytes(), content) {
		t.Errorf("downloaded = %q, want %q", buf.Bytes(), content)
	}
}

func TestDownloadURLNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	var buf bytes.Buffer
	if err := DownloadURL(context.Background(), nil, srv.URL, &buf); err == nil {
		t.Error("DownloadURL() error = nil, want error for 404")
	}
}
