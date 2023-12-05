package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies
// for the web application. For now we'll only include fields for the
// two custom loggers, but we'll add more to it as the build progresses.

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Use log.New() to create a logger for writing information messages.
	// This takes three parameters: the destination to write the logs to (os.Stdout),
	// a string prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time).
	// Note that the flags are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant file name and line number.
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Initialize a new http.Server struct.
	// We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before,
	// and set the ErrorLog field so that the server now uses the custom
	// errorLog logger in the event of any problems.
	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Call the new app.routes() method to get the servemux containing our routes.
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)

	// The value returned from the flag.String() function is a pointer to the flag value,
	// not the value itself. So we need to dereference the pointer before using it.
	err := server.ListenAndServe()
	errorLog.Fatal(err)
}
