package user

import (
	"context"
	"database/sql"
	"log"

	"github.com/TomasConesa/go-fundamentals-web-users/internal/domain"
)

type (
	Repository interface {
		Create(ctx context.Context, user *domain.User) error
		GetAll(ctx context.Context) ([]domain.User, error)
		GetById(ctx context.Context, id uint64) (*domain.User, error)
		Update(ctx context.Context, id uint64, firstName, lastName, email *string) error
	}

	repo struct {
		db  *sql.DB
		log *log.Logger
	}
)

func NewRepo(db *sql.DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (r *repo) Create(ctx context.Context, user *domain.User) error {
	/* 	r.db.MaxUserId++
	   	user.Id = r.db.MaxUserId
	   	r.db.Users = append(r.db.Users, *user) */
	return nil
}

func (r *repo) GetAll(ctx context.Context) ([]domain.User, error) {
	return nil, nil
}

func (r *repo) GetById(ctx context.Context, id uint64) (*domain.User, error) {
	/* 	index := slices.IndexFunc(r.db.Users, func(u domain.User) bool {
	   		return u.Id == id
	   	})

	   	if index < 0 {
	   		return nil, ErrNotFound{id}
	   	} */
	return nil, nil
}

func (r *repo) Update(ctx context.Context, id uint64, firstName, lastName, email *string) error {
	/* 	user, err := r.GetById(ctx, id)
	   	if err != nil {
	   		return err
	   	}

	   	if firstName != nil {
	   		user.FirstName = *firstName
	   	}

	   	if lastName != nil {
	   		user.LastName = *lastName
	   	}

	   	if email != nil {
	   		user.Email = *email
	   	} */

	return nil
}
