package datastore

type database struct {
	connection bool
	users      Userstore
}

var db Datastore

func newDefaultDatabase() Datastore {
	return &database{
		connection: true,
		users:      newUserDBStore(),
	}
}
