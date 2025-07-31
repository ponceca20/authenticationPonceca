package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword generates a bcrypt hash of the password.
// It uses a cost factor of 12, which is a good balance between security and performance.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// CheckPasswordHash compares a bcrypt hashed password with its possible plaintext equivalent.
// It returns true if the password and hash match, and false otherwise.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
