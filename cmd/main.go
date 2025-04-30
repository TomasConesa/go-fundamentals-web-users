package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TomasConesa/go-fundamentals-web-users/internal/user"
	"github.com/TomasConesa/go-fundamentals-web-users/pkg/bootstrap"
	"github.com/TomasConesa/go-fundamentals-web-users/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {

	/* 	err := godotenv.Load() // Como ya tengo un archivo .env no hace falta pasarle nada
	   	if err != nil {
	   		log.Fatal("Error loading .env file")
	   	} */

	_ = godotenv.Load()

	server := http.NewServeMux()

	db, err := bootstrap.NewDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() //Cierro la conexión a la base de datos

	if err := db.Ping(); err != nil { // Para chequear la conexipon del usuario correspondiente a la base de datos
		log.Fatal(err)
	}

	logger := bootstrap.NewLogger()
	repo := user.NewRepo(db, logger)
	service := user.NewService(logger, repo)

	ctx := context.Background()

	handler.NewUserHTTPServer(ctx, server, user.MakeEndpoints(ctx, service))

	port := os.Getenv("PORT")
	fmt.Println("Server started at port ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), server))
}
