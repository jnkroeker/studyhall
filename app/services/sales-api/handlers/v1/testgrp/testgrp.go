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
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	if n := rand.Intn(100); n%2 == 0 {
		// return errors.New("untrusted error")
		return validate.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
		// panic("testing panic")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
