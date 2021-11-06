package testgrp

import (
<<<<<<< HEAD
	"context"
	"errors"
	"jnk-ardan-service/business/sys/validate"
	"jnk-ardan-service/foundation/web"
	"math/rand"
=======
	"encoding/json"
>>>>>>> add handlers
	"net/http"

	"go.uber.org/zap"
)

// Handlers manages the set of check endpoints.
type Handlers struct {
	Log *zap.SugaredLogger
}

// Test handler is for development
<<<<<<< HEAD
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	if n := rand.Intn(100); n%2 == 0 {
		// return errors.New("untrusted error")
		return validate.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
		// panic("testing panic")
	}

=======
func (h Handlers) Test(w http.ResponseWriter, r *http.Request) {
>>>>>>> add handlers
	status := struct {
		Status string
	}{
		Status: "OK",
	}
<<<<<<< HEAD

	return web.Respond(ctx, w, status, http.StatusOK)
=======
	json.NewEncoder(w).Encode(status)

	statusCode := http.StatusOK

	h.Log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
>>>>>>> add handlers
}
