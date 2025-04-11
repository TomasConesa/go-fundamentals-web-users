package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endopoint struct {
		Create Controller
		GetAll Controller
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
)

func MakeEndpoints(ctx context.Context, s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetAllUsers(ctx, s, w)
		case http.MethodPost:
			decode := json.NewDecoder(r.Body)
			var req CreateReq
			if err := decode.Decode(&req); err != nil {
				MsgResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			PostUser(ctx, s, w, req)
		default:
			InvalidMethod(w)

		}
	}
}

func GetAllUsers(ctx context.Context, s Service, w http.ResponseWriter) {
	users, err := s.GetAll(ctx)
	if err != nil {
		MsgResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	DataResponse(w, http.StatusOK, users)
}

func PostUser(ctx context.Context, s Service, w http.ResponseWriter, data any) { // Utilizo any en lugar de interface{}. Es lo mismo
	req := data.(CreateReq)

	if req.FirstName == "" {
		MsgResponse(w, http.StatusBadRequest, "first name is required")
		return // Si no coloco el return el codigo sigue ejecutandose y se agrega de todas maneras un usuario con el campo vacio
	}
	if req.LastName == "" {
		MsgResponse(w, http.StatusBadRequest, "last name is required")
		return
	}
	if req.Email == "" {
		MsgResponse(w, http.StatusBadRequest, "email name is required")
		return
	}

	user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email)
	if err != nil {
		MsgResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	DataResponse(w, http.StatusCreated, user)
}

func InvalidMethod(w http.ResponseWriter) {
	status := http.StatusNotFound
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "method doesnt exist"}`, status)
}

func MsgResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, message)

}

func DataResponse(w http.ResponseWriter, status int, users any) {
	value, err := json.Marshal(users)
	if err != nil {
		MsgResponse(w, http.StatusBadRequest, err.Error())
		return // para que no se siga ejecutando
	}
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, value)
}
