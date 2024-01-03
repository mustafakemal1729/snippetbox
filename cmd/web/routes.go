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

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Create a file server to serve static files from the "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static"))
	// Serve static files from the "/static" path.
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Register other application routes.
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/:id", dynamic.ThenFunc(app.snippetView))

	// Use POST for creating a new snippet (resource-oriented URL).
	router.Handler(http.MethodGet, "/snippet", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet", dynamic.ThenFunc(app.snippetCreateNote))

	// Middleware chain containing 'standard' middleware.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Apply the standard middleware to the router.
	return standardMiddleware.Then(router)
}
