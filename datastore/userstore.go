package datastore

import (
	"database/sql"
	"errors"
	"log"

	"github.com/hoopra/api-base_go/models"
	"github.com/hoopra/api-base_go/utils"
	uuid "github.com/satori/go.uuid"
)

type Userstore interface {
	Add(user *models.User) error
	Validate(user *models.User) bool
	UpdateName(id uuid.UUID, newName string) error
	UpdatePassword(id uuid.UUID, newPassword string) error
	GetUUIDByName(name string) (uuid.UUID, error)
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

	existing := store.SelectByName(user.Username)
	if existing != nil {
		return errors.New("A user with that name already exists.")
	}

	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	newUser := User{
		UUID: uuid.NewV4()
		Username: user.Username
		Hash: hash
	}

	err := s.db.Exec("INSERT INTO users (id, name, hash) VALUES ($1, NULLIF($2,''), $3)",
		newUser.UUID,
		newUser.Username,
		newUser.Hash,
	)

	return nil
}

func (s *UserDBStore) Validate(user *models.User) bool {

	dbUser := store.SelectByName(user.Username)
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
		uuid,
	)

	return err
}

func (s *UserDBStore) UpdatePassword(id uuid.UUID, newPassword string) error {

	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	_, err := s.db.Exec(`UPDATE users SET password_hash = $1 WHERE id = $2`,
		hash,
		uuid,
	)

	return err
}

func (s *UserDBStore) GetUUIDByName(name string) (uuid.UUID, error) {

	user := store.SelectByName(name)
	if user == nil {
		return uuid.FromStringOrNil(""), errors.New("No such user")
	}

	return user.UUID, nil
}

func (s *UserDBStore) SelectByName(name string) *User, err {

	user := User{}
	err := s.tx.QueryRow(`SELECT (id, name, password_hash) FROM users WHERE name = $1`,
		name,
	).Scan(
		&user.UUID,
		&user.Username,
		&user.Hash,
	)

	return &user, err
}

func (s *UserDBStore) SelectByID(id uuid.UUID) *User {

	user := User{ UUID: id }
	err := s.tx.QueryRow(`SELECT (name, password_hash) FROM users WHERE id = $1`,
		id,
	).Scan(
		&user.Username,
		&user.Hash,
	)

	return &user, err
}
