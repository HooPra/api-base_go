package datastore

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/hoopra/api-base_go/models"
	"github.com/hoopra/api-base_go/utils"
	uuid "github.com/satori/go.uuid"
)

type Userstore interface {
	Add(user *models.User) error
	Validate(user *models.User) bool
	UpdateName(id uuid.UUID, newName string) error
	UpdatePassword(id uuid.UUID, newPassword string) error
	// GetUUIDByName(name string) (uuid.UUID, error)
	SelectByID(id uuid.UUID) *User
	SelectByName(name string) *User
}

type User struct {
	UUID     uuid.UUID
	Username string
	Hash     []byte
}

type UserDBStore struct {
	db *sql.DB
}

func newUserStore(db *sql.DB) *UserDBStore {
	return &UserDBStore{db: db}
}

func (s *UserDBStore) Add(user *models.User) error {

	existing := s.SelectByName(user.Username)

	if existing != nil {
		return errors.New("a user with that name already exists")
	}

	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	newUser := User{
		UUID:     id,
		Username: user.Username,
		Hash:     hash,
	}

	_, err = s.db.Exec("INSERT INTO users (id, name, password_hash) VALUES ($1, NULLIF($2,''), $3)",
		newUser.UUID,
		newUser.Username,
		newUser.Hash,
	)

	return nil
}

func (s *UserDBStore) Validate(user *models.User) bool {

	dbUser := s.SelectByName(user.Username)
	if dbUser == nil {
		return false
	}

	err := utils.CompareHashAndPassword(dbUser.Hash, user.Password)
	if dbUser.Username == user.Username && err == nil {
		return true
	}

	return false
}

func (s *UserDBStore) UpdateName(id uuid.UUID, newName string) error {

	_, err := s.db.Exec(`UPDATE users SET name = $1 WHERE id = $2`,
		newName,
		id,
	)

	return err
}

func (s *UserDBStore) UpdatePassword(id uuid.UUID, newPassword string) error {

	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`UPDATE users SET password_hash = $1 WHERE id = $2`,
		hash,
		id,
	)

	return err
}

// func (s *UserDBStore) GetUUIDByName(name string) (uuid.UUID, error) {

// 	user, err := s.SelectByName(name)
// 	if user == nil {
// 		return uuid.FromStringOrNil(""), errors.New("no such user")
// 	}

// 	return user.UUID, nil
// }

func (s *UserDBStore) SelectByName(name string) *User {

	user := User{}
	row := s.db.QueryRow(`SELECT id, name, password_hash FROM users WHERE name = $1`,
		name,
	)

	switch err := row.Scan(&user.UUID, &user.Username, &user.Hash); err {
	case sql.ErrNoRows:
		return nil
		// case nil:
		// default:
	}

	return &user
}

func (s *UserDBStore) SelectByID(id uuid.UUID) *User {

	user := User{UUID: id}
	row := s.db.QueryRow(`SELECT (name, password_hash) FROM users WHERE id = $1`,
		id,
	)

	switch err := row.Scan(&user.Username, &user.Hash); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return nil
	case nil:
	default:
		panic(err)
	}

	return &user
}
