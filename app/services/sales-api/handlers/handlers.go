// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"expvar"
	"jnk-ardan-service/app/services/sales-api/handlers/debug/checkgrp"
	"jnk-ardan-service/app/services/sales-api/handlers/v1/testgrp"
<<<<<<< HEAD
	"jnk-ardan-service/business/web/mid"
	"jnk-ardan-service/foundation/web"

=======
	"jnk-ardan-service/foundation/web"
>>>>>>> return the user-defined type web.App which embeds the http mux, abstracting the specific mux in use
	"net/http"
	"net/http/pprof"
	"os"

	"go.uber.org/zap"
)

// DebugStandardLibraryMux registers all the debug routes from the standard lilbrary
// into a new mux bypassing the use of the DefaultServerMux.
// Using the DefaultServerMux would be a security risk since a dependency could inject a
// handler into our service without us knowing it.
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

func DebugMux(build string, log *zap.SugaredLogger) http.Handler {
	mux := DebugStandardLibraryMux()

	// Register debug check endpoints.
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// APIMux constructs an http.Handler with all applications routes defined.
func APIMux(cfg APIMuxConfig) *web.App {
<<<<<<< HEAD

	// Construct the web.App "onion" which holds all routes.
	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		mid.Panics(),
	)

	// Load the routes for the different versions of the API
	v1(app, cfg)

	return app
}

func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"
=======
	app := web.NewApp(cfg.Shutdown)
>>>>>>> return the user-defined type web.App which embeds the http mux, abstracting the specific mux in use

	tgh := testgrp.Handlers{
		Log: cfg.Log,
	}
<<<<<<< HEAD

	app.Handle(http.MethodGet, "v1", "/test", tgh.Test)
=======
	app.Handle(http.MethodGet, "/v1/test", tgh.Test)

	return app
>>>>>>> return the user-defined type web.App which embeds the http mux, abstracting the specific mux in use
}
