package server

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


// NewServer


func TestNewServer_SetsAddr(t *testing.T) {
	t.Parallel()

	srv := NewServer(":9090", nil, nil)
	if srv.Addr != ":9090" {
		t.Errorf("Addr: want %q, got %q", ":9090", srv.Addr)
	}
}

func TestNewServer_NilDependencies(t *testing.T) {
	t.Parallel()

	srv := NewServer(":8080", nil, nil)
	if srv.Client != nil {
		t.Error("Client should be nil when passed nil")
	}
	if srv.Publisher != nil {
		t.Error("Publisher should be nil when passed nil")
	}
}


// RegisterRoutes – route existence checks


func TestRegisterRoutes_RootReturnsOK(t *testing.T) {
	t.Parallel()

	srv := NewServer(":0", nil, nil)
	handler := srv.RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GET /: want 200, got %d", rec.Code)
	}
}

func TestRegisterRoutes_ProjectsReturnsOK(t *testing.T) {
	t.Parallel()

	srv := NewServer(":0", nil, nil)
	handler := srv.RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GET /projects: want 200, got %d", rec.Code)
	}
}

func TestRegisterRoutes_GamesColorShooterReturnsOK(t *testing.T) {
	t.Parallel()

	srv := NewServer(":0", nil, nil)
	handler := srv.RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/games/color-shooter", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Should succeed (Publisher is nil, so it logs a warning and continues)
	if rec.Code != http.StatusOK {
		t.Errorf("GET /games/color-shooter: want 200, got %d", rec.Code)
	}
}


// handlerProjects – HTMX vs full page


func TestHandlerProjects_FullPage_ContainsHTML(t *testing.T) {
	t.Parallel()

	srv := NewServer(":0", nil, nil)
	handler := srv.RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "High-Performance Portfolio") {
		t.Error("full-page /projects response should contain project titles")
	}
}

func TestHandlerProjects_HTMXRequest_ReturnsPartial(t *testing.T) {
	t.Parallel()

	srv := NewServer(":0", nil, nil)
	handler := srv.RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("HTMX GET /projects: want 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	// HTMX partial should NOT contain the full <html> wrapper
	if strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("HTMX partial should not return full HTML document")
	}
}


// handlerUserStartGame – nil publisher path


func TestHandlerUserStartGame_NilPublisher_StillServes(t *testing.T) {
	t.Parallel()

	srv := NewServer(":0", nil, nil)
	handler := srv.RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/games/color-shooter", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200 even with nil publisher, got %d", rec.Code)
	}
}


// addGodotHeaders middleware


func TestAddGodotHeaders_SetsCOOP(t *testing.T) {
	t.Parallel()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := addGodotHeaders(inner)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	got := rec.Header().Get("Cross-Origin-Opener-Policy")
	if got != "same-origin" {
		t.Errorf("COOP header: want %q, got %q", "same-origin", got)
	}
}

func TestAddGodotHeaders_SetsCOEP(t *testing.T) {
	t.Parallel()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := addGodotHeaders(inner)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	got := rec.Header().Get("Cross-Origin-Embedder-Policy")
	if got != "require-corp" {
		t.Errorf("COEP header: want %q, got %q", "require-corp", got)
	}
}

func TestAddGodotHeaders_PassesThroughBody(t *testing.T) {
	t.Parallel()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	handler := addGodotHeaders(inner)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Body.String() != "ok" {
		t.Errorf("body: want %q, got %q", "ok", rec.Body.String())
	}
}


// gzipMiddleware


func TestGzipMiddleware_CompressesWhenAccepted(t *testing.T) {
	t.Parallel()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	handler := gzipMiddleware(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatal("expected Content-Encoding: gzip header")
	}

	reader, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to decompress: %v", err)
	}

	if string(decompressed) != "hello world" {
		t.Errorf("decompressed body: want %q, got %q", "hello world", string(decompressed))
	}
}

func TestGzipMiddleware_NoCompressionWithoutHeader(t *testing.T) {
	t.Parallel()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	handler := gzipMiddleware(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// No Accept-Encoding header
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Content-Encoding") == "gzip" {
		t.Error("should not set gzip encoding when client doesn't accept it")
	}
	if rec.Body.String() != "hello world" {
		t.Errorf("body: want %q, got %q", "hello world", rec.Body.String())
	}
}

func TestGzipMiddleware_CompressedBodyDiffersFromPlain(t *testing.T) {
	t.Parallel()

	payload := "hello world"
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(payload))
	})
	handler := gzipMiddleware(inner)

	// With gzip
	gzReq := httptest.NewRequest(http.MethodGet, "/", nil)
	gzReq.Header.Set("Accept-Encoding", "gzip")
	gzRec := httptest.NewRecorder()
	handler.ServeHTTP(gzRec, gzReq)

	// Without gzip
	plainReq := httptest.NewRequest(http.MethodGet, "/", nil)
	plainRec := httptest.NewRecorder()
	handler.ServeHTTP(plainRec, plainReq)

	// Raw bytes on the wire should differ (compressed vs plain)
	if gzRec.Body.String() == plainRec.Body.String() {
		t.Error("gzip-encoded body should differ from plain body on the wire")
	}
}

func TestGzipMiddleware_HandlesDeflateOnly(t *testing.T) {
	t.Parallel()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("plain"))
	})
	handler := gzipMiddleware(inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "deflate")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Content-Encoding") == "gzip" {
		t.Error("should not gzip when only deflate is accepted")
	}
	if rec.Body.String() != "plain" {
		t.Errorf("body: want %q, got %q", "plain", rec.Body.String())
	}
}


// gzipResponseWriter


func TestGzipResponseWriter_WritesThrough(t *testing.T) {
	t.Parallel()

	var buf strings.Builder
	rec := httptest.NewRecorder()
	grw := gzipResponseWriter{
		Writer:         &buf,
		ResponseWriter: rec,
	}

	n, err := grw.Write([]byte("test"))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if n != 4 {
		t.Errorf("bytes written: want 4, got %d", n)
	}
	if buf.String() != "test" {
		t.Errorf("writer received: want %q, got %q", "test", buf.String())
	}
	// The underlying ResponseWriter should NOT have the data
	// since Write goes to the Writer (gzip), not ResponseWriter
	if rec.Body.Len() != 0 {
		t.Error("ResponseWriter body should be empty when Writer intercepts")
	}
}
