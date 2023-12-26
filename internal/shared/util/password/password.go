package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pass, salt string) (string, error) {
	saltedPass := fmt.Sprintf("%s%s", pass, salt)
	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparedPwd(plainPass, hashedPass, salt string) error {
	saltedPass := fmt.Sprintf("%s%s", plainPass, salt)
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(saltedPass))
}
