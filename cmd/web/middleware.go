package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Open the log file in append mode
		logFile, err := os.OpenFile("access.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			// Handle the error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer logFile.Close()

		// Create a multi writer to write to log file
		multiWriter := io.MultiWriter(logFile)

		// Create a new logger that writes to the multi writer
		logger := log.New(multiWriter, "", log.LstdFlags)

		// Log the request and response details
		logger.Printf("%s -- %s %s %s %s ", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI(), r.Host)

		// Call the next handler with the ResponseRecorder instead of the original ResponseWriter
		next.ServeHTTP(w, r)

	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("panic recovered: %s", err))

			}
		}()
		next.ServeHTTP(w, r)
	})
}
