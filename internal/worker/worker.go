package worker

import (
	"context"
	"load-balancer/internal/helper/http_helper"
	"load-balancer/internal/helper/log_helper"
	"net"
	"net/http"
	"time"
)

func Run(pctx context.Context, workerId int, listenAddr string, timeout time.Duration) (err error) {
	cctx, cancel := context.WithTimeout(pctx, timeout)
	defer cancel()

	handler := newHttpHandler(cctx, workerId)
	server := http.Server{
		Handler: handler,
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log_helper.LogError(cctx, "Failed to create a listener on %s. Error: %v", listenAddr, err)
		return
	}

	err = http_helper.InterruptableServe(cctx, listener, &server, serverShutdownTime)
	if err != nil {
		log_helper.LogError(cctx, "Failed to serve on %s. Error: %v", listenAddr, err)
		return
	}

	return
}
