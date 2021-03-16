package main

import (
	"net/http"
	"time"

	"github.com/dre1080/recovr"
	"github.com/labstack/gommon/log"
	"github.com/rs/cors"
)

type captureResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newCaptureResponseWriter(w http.ResponseWriter) *captureResponseWriter {
	return &captureResponseWriter{w, http.StatusOK}
}

func (lrw *captureResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.SetHeader(`{"time":"${time_rfc3339_nano}","level":"${level}"}`)

		start := time.Now()
		lrw := newCaptureResponseWriter(w)

		next.ServeHTTP(lrw, r)

		elapsed := time.Since(start)
		code := lrw.statusCode
		if code >= 500 {
			log.Errorf("[ACCESS] %v %v %v %v\n", r.Method, code, r.URL, elapsed)
		} else if code >= 400 {
			log.Warnf("[ACCESS] %v %v %v %v\n", r.Method, code, r.URL, elapsed)
		} else {
			log.Infof("[ACCESS] %v %v %v %v\n", r.Method, code, r.URL, elapsed)
		}
	})
}

// @see panic-handling https://github.com/go-swagger/go-swagger/blob/master/docs/use/middleware.md#add-logging-and-panic-handling
func recoveryMiddleware(next http.Handler) http.Handler {
	recovery := recovr.New()
	return recovery(next)
}

// @see https://github.com/go-swagger/go-swagger/blob/master/docs/faq/faq_documenting.md#how-to-use-swagger-ui-cors
func corsMiddleware(next http.Handler) http.Handler {
	handleCORS := cors.AllowAll().Handler
	return handleCORS(next)
}
