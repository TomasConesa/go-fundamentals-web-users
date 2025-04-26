package user

import (
	"context"
	"errors"

	"github.com/TomasConesa/go-fundamentals-response/response"
)

type (
	Controller func(ctx context.Context, request any) (any, error)

	Endpoint struct {
		Create  Controller
		GetAll  Controller
		GetById Controller
		Update  Controller
	}

	GetReq struct {
		Id uint64
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	UpdateReq struct {
		Id        uint64
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
	}
)

func MakeEndpoints(ctx context.Context, s Service) Endpoint {
	return Endpoint{
		Create:  makeCreateEndpoint(s),
		GetAll:  makeGetAllEndpoint(s),
		GetById: makeGetByIdEndpoint(s),
		Update:  makeUpdateEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(CreateReq)

		if req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}
		if req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}
		if req.Email == "" {
			return nil, response.BadRequest(ErrEmailRequired.Error())
		}

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		return response.Created("success", user), nil
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {
		users, err := s.GetAll(ctx)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", users), nil
	}
}

func makeGetByIdEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetReq)

		user, err := s.GetById(ctx, req.Id)
		if err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", user), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(UpdateReq)

		if req.FirstName != nil && *req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}
		if req.LastName != nil && *req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}
		if req.Email != nil && *req.Email == "" {
			return nil, response.BadRequest(ErrEmailRequired.Error())
		}
		if err := s.Update(ctx, req.Id, req.FirstName, req.LastName, req.Email); err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", nil), nil
	}
}
