package repository

import (
	"fmt"
	"net/mail"

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
	db *database.Service,
) UserRepository {
	return UserRepository{
		db: *db,
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

func (s *UserRepository) GetUserByEmail(email string) (*database.User, error) {
	row := s.db.QueryRow(
		"SELECT id, username, first_name, last_name, email, password FROM users WHERE email = $1;",
		email,
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

// Validate User Fields
func (s *UserRepository) ValidateUserFields(user database.User) error {
	if user.Username == "" {
		return fmt.Errorf("username is required")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return fmt.Errorf("invalid email address")
	}

	if user.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if user.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if len(user.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	return nil
}
