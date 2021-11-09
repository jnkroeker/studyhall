// Package web contains a small web framework extension
package web

import (
<<<<<<< HEAD
	"context"
	"net/http"
=======
>>>>>>> return the user-defined type web.App which embeds the http mux, abstracting the specific mux in use
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux/v5"
)

<<<<<<< HEAD
// A Handler is a type that handles an http request within our mini framework
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

=======
>>>>>>> return the user-defined type web.App which embeds the http mux, abstracting the specific mux in use
// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
<<<<<<< HEAD
	mw       []Middleware
}

// NewApp creates an App value that handles a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
=======
}

// NewApp creates an App value that handles a set of routes for the application.
func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
>>>>>>> return the user-defined type web.App which embeds the http mux, abstracting the specific mux in use
	}
}

// SginalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}
<<<<<<< HEAD

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) Handle(method string, group string, path string, handler Handler, mw ...Middleware) {

	// First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(mw, handler)

	// add the application's general middleware to the handler chain
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {

		// INJECT CODE

		// Call the wrapped handler functions.
		if err := handler(r.Context(), w, r); err != nil {

			// INJECT CODE
			return
		}

		// INJECT CODE
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	a.ContextMux.Handle(method, finalPath, h)
}
=======
>>>>>>> return the user-defined type web.App which embeds the http mux, abstracting the specific mux in use
