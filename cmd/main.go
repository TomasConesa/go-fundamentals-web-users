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

	/* server := http.NewServeMux() */ // Nativo de go

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

	/* handler.NewUserHTTPServer(ctx, server, user.MakeEndpoints(ctx, service)) */ // Con la forma nativa de go
	h := handler.NewUserHTTPServer(user.MakeEndpoints(ctx, service))

	port := os.Getenv("PORT")
	fmt.Println("Server started at port ", port)

	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Handler: accessControl(h),
		Addr:    address,
	}

	log.Fatal(srv.ListenAndServe())
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, HEAD, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With,X-Forwarded-For")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
