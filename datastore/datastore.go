package datastore

import (
	"database/sql"
)

type Database interface {
	Rollback()
}

struct Datastore {
	connection *sql.DB
	users      *Userstore
}

var datastore Datastore

func getDatastore() *Datastore {
	if datastore == nil {
		datastore = newDefaultDatabase()
	}
	return datastore
}

func newDefaultDatastore() *Datastore {
	conn, err := connect()
	if err != nil {
		panic("cannot connect to db")
	}
	datastore = {
		connection: conn
		users: newUserStore(conn)
	}
}

// Connection returns a connection to the underlying database
func connect() *sql.DB {
	// TODO read credentials from file
	// TODO accept multiple driver types
	// conn, err := sql.Open("postgres", "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")
	conn, err := sql.Open("mysql", "user:password@/authapi")
	if err != nil {
		panic("could not connect to db")
	}
	return conn
}

/* Native functionality */
func (db *Datastore) Rollback() {}
