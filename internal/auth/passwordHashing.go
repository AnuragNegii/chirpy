package auth 

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)



func HashPassword(password string)(string, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", fmt.Errorf("password couln't be hashed: %v", err)
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error{
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
