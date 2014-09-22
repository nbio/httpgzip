package httpgzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	gzipWriter io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.gzipWriter.Write(b)
}

// GzipResponse wraps an http.Handler and compresses the response
// if requested in the Accept-Encoding header.
func GzipResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Vary", "Accept-Encoding")
		if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, req)
			return
		}
		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()
		gw := &gzipResponseWriter{
			ResponseWriter: w,
			gzipWriter: gzipWriter,
		}
		gw.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gw, req)
	})
}
