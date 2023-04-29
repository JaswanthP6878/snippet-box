package main

import (
	"fmt"
	"net/http"
)

// func secureHeaders(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
// 		w.Header().Set("Refferer-Policy", "origin-when-cross-origin")
// 		w.Header().Set("X-Content-Type-Options", "nosniff")
// 		w.Header().Set("X-Frame-Options", "deny")
// 		w.Header().Set("X-XSS-Protection", "0")

// 		next.ServeHTTP(w, r)
// 	})
// }

func (app *application) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s- %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// this defered call gets executed in the unwinding stage
		// if a panic occurs
		defer func() {
			if err := recover(); err != nil {
				// sets a "Connection: close" header on the response
				w.Header().Set("Connection", "close")
				// Internal ServerError
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)

	})
}