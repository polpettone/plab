package pkg

import "golang.org/x/crypto/bcrypt"

func bcryptHash(value string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), 14)
	return string(bytes), err
}

