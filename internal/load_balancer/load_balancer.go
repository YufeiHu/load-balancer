package load_balancer

import (
	"context"
	"load-balancer/internal/helper/context_helper"
	"load-balancer/internal/helper/log_helper"
	"load-balancer/internal/helper/uuid_helper"
	"net"
	"net/http"
	"sync"
)

type worker struct {
	key    string
	conn   net.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

func (w *worker) Close() {
	if w.cancel != nil {
		w.cancel()
	}

	if w.conn != nil {
		_ = w.conn.Close()
	}
}

type loadBalancer struct {
	workerLock sync.Mutex
	workerCond *sync.Cond
	workerSet  map[string]*worker
}

func (lb *loadBalancer) wakeWaiters() {
	lb.workerCond.Broadcast()
}

func (lb *loadBalancer) monitorWorkers(ctx context.Context, listenerOnWorkers net.Listener) error {
	defer func() {
		_ = listenerOnWorkers.Close()
	}()

	for ctx.Err() == nil {
		conn, err := listenerOnWorkers.Accept()
		if err != nil {
			log_helper.LogError(ctx, "Failed to accept requests from workers. Error: %v", err)
			return err
		}

		log_helper.LogInfo(ctx, "Adding worker: %v", conn.RemoteAddr())
		// TODO
		// lb.addWorker(ctx, conn)
	}

	return ctx.Err()
}

func (lb *loadBalancer) handleClientRequests(pctx context.Context, listenerOnClient net.Listener) error {
	for pctx.Err() == nil {
		conn, err := listenerOnClient.Accept()
		log_helper.LogInfo(pctx, "Received a connection object from: %v", conn.RemoteAddr())
		if err != nil {
			log_helper.LogError(pctx, "Failed to handle a client request. Error: %v", err)
			continue
		}

		uuid := uuid_helper.GenerateUuid()
		cctx := context_helper.SetWorkRequestId(pctx, uuid)

		//handler := newHttpHandler(cctx)
		err = http.Serve(connListener{conn, make(chan struct{})}, nil)
		if err != nil {
			log_helper.LogError(cctx, "Failed to handle a client request. Error: %v", err)
			continue
		}

		//conn.Write([]byte("Hello"))
		log_helper.LogInfo(cctx, "Successfully handled a client request")
	}

	return pctx.Err()
}

func Run(ctx context.Context, listenerOnClient net.Listener) error {
	lb := &loadBalancer{
		workerSet: make(map[string]*worker),
	}

	//handler := newHttpHandler(ctx)
	//server := http.Server{
	//	Handler: handler,
	//}
	//
	//server.Serve()

	lb.workerCond = sync.NewCond(&lb.workerLock)

	go func() {
		defer func() {
			_ = listenerOnClient.Close()
			lb.wakeWaiters()
		}()

		<-ctx.Done()
		log_helper.LogWarn(ctx, "Shutting down load balancer")
	}()

	//go func() {
	//	var listenerOnWorkers net.Listener
	//	err := lb.monitorWorkers(ctx, listenerOnWorkers)
	//	if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
	//		log_helper.LogError(ctx, "Failed to monitor workers. Error is: %v", err)
	//	}
	//}()

	return lb.handleClientRequests(ctx, listenerOnClient)
}
