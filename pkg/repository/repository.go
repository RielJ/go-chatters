package repository

import "github.com/rielj/go-chatters/pkg/database"

type Repository struct {
	User UserRepository
}

func New(db *database.Service) Repository {
	return Repository{
		User: NewUserRepository(db),
	}
}
