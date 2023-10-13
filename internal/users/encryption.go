package users

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (string, error) {
	if password == "" {
		return "", nil
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CompareHashAndPassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
