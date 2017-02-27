package datastore

import "golang.org/x/crypto/bcrypt"

// HashPassword returns a hash of a password
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 10)
}

func CompareHashAndPassword(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
