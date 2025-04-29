package user

import (
	"errors"
	"fmt"
)

var ErrFirstNameRequired = errors.New("first name is required")
var ErrLastNameRequired = errors.New("last name is required")
var ErrEmailRequired = errors.New("email is required")
var ErrThereArentFields = errors.New("there aren't fields")

type ErrNotFound struct {
	Id uint64
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user id '%d' doesn't exist", e.Id)
}
