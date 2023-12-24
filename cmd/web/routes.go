package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// routes configures and returns the application's HTTP router.
func (app *application) routes() http.Handler {
	// Create a new router instance from httprouter.
	router := httprouter.New()

	// Create a file server to serve static files from the "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Serve static files from the "/static" path.
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Register other application routes.
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/:id", app.snippetView)

	// Use POST for creating a new snippet (resource-oriented URL).
	router.HandlerFunc(http.MethodPost, "/snippets", app.snippetCreate)

	// Middleware chain containing 'standard' middleware.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Apply the standard middleware to the router.
	return standardMiddleware.Then(router)
}
