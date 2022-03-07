package worker

import (
	"context"
	"encoding/json"
	"load-balancer/internal/helper/context_helper"
	"load-balancer/internal/helper/http_helper"
	"net/http"
)

type response struct {
	WorkerId      int `json:"workerId"`
	WorkRequestId string `json:"workRequestId"`
}

func newHttpHandler(ctx context.Context, workerId int) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(basePath + "/", func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			_ = req.Body.Close()
		}()

		if req.Method != "GET" {
			http_helper.WriteStatus(ctx, w, req, http.StatusMethodNotAllowed, nil)
			return
		}

		r := response{
			WorkerId: workerId,
			WorkRequestId: context_helper.WorkRequestId(ctx),
		}

		body, err := json.Marshal(r)
		if err != nil {
			http_helper.WriteStatus(ctx, w, req, http.StatusInternalServerError, err)
			return
		}

		http_helper.WriteStatus(ctx, w, req, http.StatusOK, err)
		_, err = w.Write(body)
		return
	})

	return mux
}
