package passwordhelpers

// PasswordHasher is an interface for hashing and comparing passwords
type PasswordHasher interface {
	// HashPassword hashes the passwords and returns the hash. Also returns any errors.
	HashPassword(password string) ([]byte, error)

	// ComparePasswords compares a password hash and a plain text password and returns any errors.
	ComparePasswords(hash []byte, password string) error
}
