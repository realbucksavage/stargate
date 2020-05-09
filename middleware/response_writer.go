package middleware

import "net/http"

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader makes loggingResponseWriter implement http.ResponseWriter.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}
