package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TomasConesa/go-fundamentals-web-users/internal/domain"
	"github.com/TomasConesa/go-fundamentals-web-users/internal/user"
)

func main() {
	server := http.NewServeMux()

	db := user.DB{
		Users: []domain.User{{
			Id:        1,
			FirstName: "Tomas",
			LastName:  "Conesa",
			Email:     "toto@gmail.com",
		}, {
			Id:        2,
			FirstName: "Toteiro",
			LastName:  "00",
			Email:     "toteir00@gmail.com",
		}, {
			Id:        3,
			FirstName: "Pony",
			LastName:  "Salvage",
			Email:     "pony@gmail.com",
		}},
		MaxUserId: 3,
	}

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	repo := user.NewRepo(db, logger)
	service := user.NewService(logger, repo)

	ctx := context.Background()

	server.HandleFunc("/users", user.MakeEndpoints(ctx, service))
	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}
