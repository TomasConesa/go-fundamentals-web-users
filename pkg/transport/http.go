package transport

import (
	"context"
	"net/http"
	"strings"
)

type Endpoint func(ctx context.Context, request any) (any, error)

type Transport interface {
	Server(
		endpoint Endpoint,
		decode func(ctx context.Context, r *http.Request) (any, error), // Actua en el medio de la peticion y el controlador
		encode func(ctx context.Context, w http.ResponseWriter, response any) error, // Actua entre el controlador y la respuesta al cliente
		encodeError func(ctx context.Context, err error, w http.ResponseWriter),
	)
}

type transport struct {
	w   http.ResponseWriter
	r   *http.Request
	ctx context.Context
}

func New(w http.ResponseWriter, r *http.Request, ctx context.Context) Transport {
	return &transport{
		w:   w,
		r:   r,
		ctx: ctx,
	}
}

func (t *transport) Server(
	endpoint Endpoint,
	decode func(ctx context.Context, r *http.Request) (any, error),
	encode func(ctx context.Context, w http.ResponseWriter, response any) error,
	encodeError func(ctx context.Context, err error, w http.ResponseWriter),
) {
	data, err := decode(t.ctx, t.r)
	if err != nil {
		encodeError(t.ctx, err, t.w)
		return
	}

	response, err := endpoint(t.ctx, data)
	if err != nil {
		encodeError(t.ctx, err, t.w)
		return
	}

	if err := encode(t.ctx, t.w, response); err != nil {
		encodeError(t.ctx, err, t.w)
		return
	}
}

func CleanUrl(url string) ([]string, int) {
	if url[0] != '/' { // Slash / al inicio de la url
		url = "/" + url
	}
	if url[len(url)-1] != '/' { // Slash / al final de la url
		url += "/"
	}

	partsUrl := strings.Split(url, "/")

	return partsUrl, len(partsUrl)
}
