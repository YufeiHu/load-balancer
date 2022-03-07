package load_balancer

import (
	"context"
	"encoding/json"
	"load-balancer/internal/config"
	"load-balancer/internal/helper/http_helper"
	"net/http"
)

type response struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func newHttpHandler(ctx context.Context) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(config.LbBasePath + "/", func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			_ = req.Body.Close()
		}()

		if req.Method != "GET" {
			http_helper.WriteStatus(ctx, w, req, http.StatusMethodNotAllowed, nil)
			return
		}

		r := response{
			"1",
			"yufhu",
		}

		body, err := json.Marshal(r)
		if err != nil {
			http_helper.WriteStatus(ctx, w, req, http.StatusInternalServerError, err)
			return
		}

		_, err = w.Write(body)
		http_helper.WriteStatus(ctx, w, req, http.StatusInternalServerError, err)
	})

	return mux
}
