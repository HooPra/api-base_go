package datastore

func Users() Userstore {
	if userstore == nil {
		db := getDatastore()
		userstore = db.users
	}
	return userstore
}
