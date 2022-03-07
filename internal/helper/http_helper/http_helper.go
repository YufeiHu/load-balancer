package http_helper

import (
	"context"
	"load-balancer/internal/helper/log_helper"
	"net"
	"net/http"
	"time"
)

func WriteStatus(ctx context.Context, w http.ResponseWriter, req *http.Request, status int, err error) {
	if err != nil {
		log_helper.LogError(ctx, "Http handler: %s %s (%s) -> %d. Error: %v", req.Method, req.RequestURI, req.RemoteAddr, status, err)
	} else {
		log_helper.LogInfo(ctx, "Http handler: %s %s (%s) -> %d", req.Method, req.RequestURI, req.RemoteAddr, status)
	}
	w.WriteHeader(status)
}

func InterruptableServe(pctx context.Context, listener net.Listener, server *http.Server, serverShutdownTime time.Duration) error {
	go func() {
		<-pctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), serverShutdownTime)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log_helper.LogError(pctx, "Failed to shutdown the server: %v", err)
			return
		}
	}()

	return server.Serve(listener)
}
