package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Id        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

var users []User

var maxId uint64

func init() {
	users = []User{{
		Id:        1,
		FirstName: "Tomas",
		LastName:  "Conesa",
		Email:     "toteiro@gmail.com",
	}, {
		Id:        2,
		FirstName: "Ruso",
		LastName:  "Ascacibar",
		Email:     "ruso@gmail.com",
	}, {
		Id:        3,
		FirstName: "Toteiro",
		LastName:  "Ascacibar",
		Email:     "sarasa@gmail.com",
	}}
	maxId = 3

}

func main() {

	http.HandleFunc("/users", UserServer)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func UserServer(w http.ResponseWriter, r *http.Request) {
	var status int
	switch r.Method {
	case http.MethodGet:
		GetAllUsers(w)
	case http.MethodPost:
		decode := json.NewDecoder(r.Body)
		var u User
		if err := decode.Decode(&u); err != nil {
			MsgResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		PostUser(w, u)
	default:
		status = 404
		w.WriteHeader(status)
		fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, "not found")

	}
}

func GetAllUsers(w http.ResponseWriter) {
	DataResponse(w, http.StatusOK, users)
}

func PostUser(w http.ResponseWriter, data interface{}) {
	user := data.(User) // Casteo la interface a tipo User
	maxId++
	user.Id = maxId
	users = append(users, user)
	DataResponse(w, http.StatusCreated, user)
}

func MsgResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, message)

}

func DataResponse(w http.ResponseWriter, status int, users interface{}) {
	value, err := json.Marshal(users)
	if err != nil {
		MsgResponse(w, http.StatusBadRequest, err.Error())
		return // para que no se siga ejecutando
	}
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, value)
}
