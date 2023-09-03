package hash

import (
	"golang.org/x/crypto/bcrypt"
)

//function to be used to hash a password
func HashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return bytes, err
}

//function to be used to compare a password with a hash
func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}