package httpgzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	gzw io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.gzw.Write(b)
}

// GzipResponse wraps an http.Handler and compresses the response
// if requested in the Accept-Encoding header.
func GzipResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, req)
			return
		}
		gzw := gzip.NewWriter(w)
		defer gzw.Close()
		gzrw := &gzipResponseWriter{gzw: gzw, ResponseWriter: w}
		gzrw.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzrw, req)
	})
}
