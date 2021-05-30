package passwordhelpers

import (
	"authserver/common"

	"golang.org/x/crypto/bcrypt"
)

// BCryptPasswordHasher is an implementation of the PasswordHasher that uses the bcrypt algorithm.
type BCryptPasswordHasher struct{}

// HashPassword hashes the password using the bcrypt algorithm and returns the hash. Also returns any errors.
func (BCryptPasswordHasher) HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, common.ChainError("bcrypt generate hash from password error", err)
	}

	return hash, nil
}

// ComparePasswords compares a password hash and a plain text password using the bcrypt algorithm and returns any errors.
func (BCryptPasswordHasher) ComparePasswords(hash []byte, password string) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return common.ChainError("bcrypt compare hash and password error", err)
	}

	return nil
}
