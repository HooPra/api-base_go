package datastore

func Users() Userstore {
	db := getDatastore()
	return db.users
}
