package models

import uuid "github.com/satori/go.uuid"

type User struct {
	UUID     uuid.UUID `json:"uuid" form:"-"`
	Username string    `json:"username" form:"username"`
	Password string    `json:"password" form:"password"`
}
