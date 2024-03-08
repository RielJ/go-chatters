package repository

import "github.com/rielj/go-chatters/pkg/database"

type Repository struct {
	User UserRepository
}

func Init(db *database.Service) Repository {
	return Repository{
		User: NewUserRepository(db),
	}
}
