package server

import (
	"compress/gzip"
	"strings"
	"io"
	"net/http"
)

type gzipResponseWriter struct {
    io.Writer
    http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}

func gzipMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }
        gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
        if err != nil {
            next.ServeHTTP(w, r)
            return
        }
        defer gz.Close()
        w.Header().Set("Content-Encoding", "gzip")
        w.Header().Del("Content-Length") // length is now unknown
        next.ServeHTTP(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
    })
}
