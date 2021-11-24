package mid

import (
	"context"
	"net/http"
	"time"

	"jnk-ardan-service/foundation/web"

	"go.uber.org/zap"
)

func Logger(log *zap.SugaredLogger) web.Middleware {

	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			traceId := "00000000000000000"
			statuscode := http.StatusOK
			now := time.Now()

			log.Infow("request started", "traceid", traceId, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Infow("request completed", "traceid", traceId, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr,
				"statuscode", statuscode, "since", time.Since(now))

			return err

		}

		return h
	}

	return m
}
