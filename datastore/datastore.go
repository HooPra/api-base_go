package datastore

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Database interface {
	Rollback()
}

type Datastore struct {
	connection *sql.DB
	users      *UserDBStore
}

var (
	retries   = 0
	datastore *Datastore
)

func Init() {
	newDefaultDatastore()
}

func getDatastore() *Datastore {
	if datastore == nil {
		datastore = newDefaultDatastore()
	}
	return datastore
}

func newDefaultDatastore() *Datastore {
	conn, _ := connect()
	// if err != nil {
	// 	panic("cannot connect to db")
	// }
	datastore = &Datastore{
		connection: conn,
		users:      newUserStore(conn),
	}
	datastore.Migrate()
	return datastore
}

// Connection returns a connection to the underlying database
func connect() (*sql.DB, error) {
	// TODO read credentials from file
	// TODO accept multiple driver types

	// sslmode=verify-full
	conn, err := sql.Open("postgres", "postgres://postgres:abc123@psql/postgres?sslmode=disable")
	// conn, err := sql.Open("mysql", "user:password@/authapi")
	if err != nil {
		if retries < 3 {
			retries++
			log.Println("retrying", retries)
			timer1 := time.NewTimer(15 * time.Second)
			<-timer1.C
			return connect()
		} else {
			panic(err)
		}
	}
	log.Println(conn, err)
	return conn, nil
}

/* Native functionality */
func (db *Datastore) Rollback() {}

func (db *Datastore) Migrate() {
	log.Println("migrating db")
	conn := db.connection
	_, err := conn.Exec(`DROP TABLE IF EXISTS users`)
	if err != nil {
		log.Println(err)
	}
	_, err = conn.Exec(`CREATE TABLE users (
			id VARCHAR PRIMARY KEY, 
			name VARCHAR (50) UNIQUE NOT NULL,
			password_hash VARCHAR NOT NULL
		)`)
	if err != nil {
		log.Println(err)
	}
}
