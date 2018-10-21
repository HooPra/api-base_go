package datastore

import (
	"errors"

	"github.com/hoopra/api-base_go/models"
	uuid "github.com/satori/go.uuid"
)

type UserDBStore map[uuid.UUID]*User

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

func newUserDBStore() UserDBStore {
	return make(map[uuid.UUID]*User)
}

func (db database) Users() Userstore {
	return db.users
}

func (store UserDBStore) Add(user *models.User) error {

	existing := store.SelectByName(user.Username)
	if existing != nil {
		return errors.New("A user with that name already exists.")
	}

	newUser := User{}
	id := uuid.NewV4()
	newUser.UUID = id
	newUser.Username = user.Username
	hash, err := HashPassword(user.Password)

	if err != nil {
		return err
	}

	newUser.Hash = hash
	store[newUser.UUID] = &newUser

	return nil
}

func (store UserDBStore) Validate(user *models.User) bool {

	dbUser := store.SelectByName(user.Username)
	if dbUser == nil {
		return false
	}

	err := CompareHashAndPassword(dbUser.Hash, user.Password)
	if dbUser.Username == user.Username && err == nil {
		return true
	}

	return false
}

func (store UserDBStore) UpdateName(id uuid.UUID, newName string) error {

	user := store.SelectByID(id)
	if user == nil {
		return errors.New("No such user")
	}

	user.Username = newName
	return nil
}

func (store UserDBStore) UpdatePassword(id uuid.UUID, newPassword string) error {

	user := store.SelectByID(id)
	if user == nil {
		return errors.New("No such user")
	}

	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Hash = hash

	return nil
}

func (store UserDBStore) GetUUIDByName(name string) (uuid.UUID, error) {

	user := store.SelectByName(name)
	if user == nil {
		return uuid.FromStringOrNil(""), errors.New("No such user")
	}

	return user.UUID, nil
}

func (store UserDBStore) SelectByName(name string) *User {

	for _, user := range store {

		if user.Username == name {
			return user
		}
	}

	return nil
}

func (store UserDBStore) SelectByID(id uuid.UUID) *User {

	return store[id]
}
