package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TomasConesa/go-fundamentals-web-users/internal/user"
	"github.com/TomasConesa/go-fundamentals-web-users/pkg/transport"
)

func NewUserHTTPServer(ctx context.Context, router *http.ServeMux, endpoint user.Endpoint) {
	router.HandleFunc("/users", UserServer(ctx, endpoint))
}

func UserServer(ctx context.Context, endpoint user.Endpoint) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tran := transport.New(w, r, ctx)

		switch r.Method {
		case http.MethodGet:
			tran.Server(
				transport.Endpoint(endpoint.GetAll),
				decodeGetAllUser,
				encodeResponse,
				encodeError,
			)
			return
		case http.MethodPost:
			tran.Server(
				transport.Endpoint(endpoint.Create),
				decodeCreateUser,
				encodeResponse,
				encodeError,
			)
			return
		}
		InvalidMethod(w)
	}
}

func decodeGetAllUser(ctx context.Context, r *http.Request) (any, error) {
	return nil, nil
}

func decodeCreateUser(ctx context.Context, r *http.Request) (any, error) {
	var req user.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request format: '%v'", err.Error())
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}
	status := http.StatusOK
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, data)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	status := http.StatusInternalServerError
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, err.Error())
}

func InvalidMethod(w http.ResponseWriter) {
	status := http.StatusNotFound
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "method doesnt exist"}`, status)
}
