package bootstrap

import (
	"log"
	"os"

	"github.com/TomasConesa/go-fundamentals-web-users/internal/domain"
	"github.com/TomasConesa/go-fundamentals-web-users/internal/user"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func NewDb() user.DB {

	return user.DB{
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
}
