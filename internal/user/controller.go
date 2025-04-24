package user

import (
	"context"
	"errors"
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
			return nil, errors.New("first name is required")
		}
		if req.LastName == "" {
			return nil, errors.New("last name is required")
		}
		if req.Email == "" {
			return nil, errors.New("email is required")
		}

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {
		users, err := s.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		return users, nil
	}
}

func makeGetByIdEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetReq)

		user, err := s.GetById(ctx, req.Id)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(UpdateReq)

		if err := s.Update(ctx, req.Id, req.FirstName, req.LastName, req.Email); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
