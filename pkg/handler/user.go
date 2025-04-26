package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/TomasConesa/go-fundamentals-response/response"
	"github.com/TomasConesa/go-fundamentals-web-users/internal/user"
	"github.com/TomasConesa/go-fundamentals-web-users/pkg/transport"
)

type contextKey string

const paramsKey contextKey = "params"

func NewUserHTTPServer(ctx context.Context, router *http.ServeMux, endpoint user.Endpoint) {
	router.HandleFunc("/users/", UserServer(ctx, endpoint))
}

func UserServer(ctx context.Context, endpoint user.Endpoint) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.Path
		log.Println(r.Method, ": ", url)

		path, pathSize := transport.CleanUrl(url)
		if pathSize < 3 || pathSize > 4 {
			InvalidMethod(w)
			return
		}

		params := make(map[string]string)
		if pathSize == 4 && path[2] != "" { // el indice 2 del path equivale al id.   /users/1 -> lo vacio es indice 0, users ind 1 y 1 el ind 2
			params["userId"] = path[2]
		}
		ctx = context.WithValue(ctx, paramsKey, params) // Agrego la clave y el valor params al contexto. Guarde la clave "params" en una constante para no pasarle el string directamente acá porque es una mala práctica

		tran := transport.New(w, r, ctx)

		var end user.Controller
		var deco func(ctx context.Context, r *http.Request) (any, error)

		switch r.Method {
		case http.MethodGet:
			switch pathSize {
			case 3:
				end = endpoint.GetAll
				deco = decodeGetAllUser
			case 4:
				end = endpoint.GetById
				deco = decodeGetUserById
			}
		case http.MethodPost:
			switch pathSize {
			case 3:
				end = endpoint.Create
				deco = decodeCreateUser
			}
		case http.MethodPatch:
			switch pathSize {
			case 4:
				end = endpoint.Update
				deco = decodeUpdateUser
			}
		}
		if end != nil && deco != nil {
			tran.Server(
				transport.Endpoint(end),
				deco,
				encodeResponse,
				encodeError,
			)
		} else {
			InvalidMethod(w)
		}
	}
}

func decodeGetAllUser(ctx context.Context, r *http.Request) (any, error) {
	return nil, nil
}

func decodeGetUserById(ctx context.Context, r *http.Request) (any, error) {
	params := ctx.Value(paramsKey).(map[string]string) // Para obtener el valor la clave params. Parseo con el map porque el metodo Value devuelve un any
	id, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		return nil, err
	}

	return user.GetReq{
		Id: id,
	}, nil
}

func decodeCreateUser(ctx context.Context, r *http.Request) (any, error) {
	var req user.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request format: '%v'", err.Error())
	}

	return req, nil
}

func decodeUpdateUser(ctx context.Context, r *http.Request) (any, error) {
	var req user.UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request format: '%v'", err.Error())
	}

	params := ctx.Value(paramsKey).(map[string]string)
	id, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		return nil, err
	}
	req.Id = id

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp any) error {

	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())

	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)

	w.WriteHeader(resp.StatusCode())

	_ = json.NewEncoder(w).Encode(resp) //  _ = ignora cualquier error que pueda pasar durante la codificación (por ejemplo, si fallara en generar el JSON).
}

func InvalidMethod(w http.ResponseWriter) {
	status := http.StatusNotFound
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "method doesn't exist"}`, status)
}
