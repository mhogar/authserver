package helpers

import (
	commonhelpers "authserver/helpers/common"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher is an interface for hashing and comparing passwords
type PasswordHasher interface {
	// HashPassword hashes the passwords and returns the hash. Also returns any errors.
	HashPassword(password string) ([]byte, error)

	// ComparePasswords compares a password hash and a plain text password and returns any errors.
	ComparePasswords(hash []byte, password string) error
}

// BCryptPasswordHasher is an implementation of the PasswordHasher that uses the bcrypt algorithm.
type BCryptPasswordHasher struct{}

// HashPassword hashes the password using the bcrypt algorithm and returns the hash. Also returns any errors.
func (BCryptPasswordHasher) HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, commonhelpers.ChainError("bcrypt generate hash from password error", err)
	}

	return hash, nil
}

// ComparePasswords compares a password hash and a plain text password using the bcrypt algorithm and returns any errors.
func (BCryptPasswordHasher) ComparePasswords(hash []byte, password string) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return commonhelpers.ChainError("bcrypt compare hash and password error", err)
	}

	return nil
}
