package datastore

import (
	"database/sql"
	"fmt"
)

var (
	connection *sql.DB
	connStr    = fmt.Sprintf("postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")
)

// type Datastore interface {
// 	Users() *sql.DB
// }

// var db struct {
// 	connection
// }

func Users() Userstore {

	if userstore == nil {
		userstore = newUserstore(Connection())
	}
	return userstore
}

// Connection returns a connection to the underlying database
func Connection() *sql.DB {
	if connection == nil {
		conn, err := sql.Open("postgres", connStr)
		if err != nil {
			panic("could not connect to db")
		}
		connection = conn
	}
	return connection
}
