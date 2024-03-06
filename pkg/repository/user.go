package repository

import (
	"fmt"

	"github.com/rielj/go-chatters/pkg/database"
	"github.com/rielj/go-chatters/pkg/tools"
)

type UserRepository struct {
	db database.Service
}

type UserRepositoryParams struct {
	Database database.Service
}

func NewUserRepository(
	parameters UserRepositoryParams,
) UserRepository {
	return UserRepository{
		db: parameters.Database,
	}
}

func (s *UserRepository) GetUserByUsername(username string) (*database.User, error) {
	row := s.db.QueryRow(
		"SELECT id, username, first_name, last_name, email, password FROM users WHERE username = $1;",
		username,
	)
	var user database.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser function
func (s *UserRepository) CreateUser(
	user database.User,
) (int64, error) {
	// hash password
	hashedPw, err := tools.HashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	fmt.Println(hashedPw, user)
	var id int64
	err = s.db.QueryRow(
		"INSERT INTO users (username, first_name, last_name, email, password) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		hashedPw,
	).Scan(&id)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return 0, err
	}
	return id, nil
}

// GetUsers function
func (s *UserRepository) GetUsers() ([]database.User, error) {
	rows, err := s.db.Query("SELECT id, username, first_name, last_name, email FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []database.User{}
	for rows.Next() {
		var user database.User
		err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUserByID function
func (s *UserRepository) GetUserByID(id string) (*database.User, error) {
	row := s.db.QueryRow(
		"SELECT id, username, first_name, last_name, email FROM users WHERE id = $1;",
		id,
	)
	var user database.User
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser function
func (s *UserRepository) UpdateUser(
	user database.User,
) error {
	_, err := s.db.Exec(
		"UPDATE users SET username = $1, first_name = $2, last_name = $3, email = $4 WHERE id = $5;",
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser function
func (s *UserRepository) DeleteUser(id string) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id = $1;", id)
	if err != nil {
		return err
	}
	return nil
}
