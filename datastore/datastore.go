package datastore

type Datastore interface {
	Users() Userstore
}

func Store() Datastore {

	if db == nil {

		// TODO: Migrate to postgres
		db = newDefaultDatabase()
	}

	return db
}
