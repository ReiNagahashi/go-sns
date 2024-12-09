package utils

import (
	"errors"
	"unicode"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string)(string, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}

	return string(hashedPassword), nil
}


func CheckPasswordHash(password, hash string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}


func ValidatePassword(password string)error{
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var hasLetter, hasNumber, hasSpecialChar bool
	for _, char := range password {
		switch {
		case unicode.IsLetter(char):
			hasLetter = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecialChar = true
		case unicode.IsSpace(char):
			return errors.New("password must not contain spaces")
		}
	}

	if !hasLetter {
		return errors.New("password must contain at least one letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecialChar {
		return errors.New("password must contain at least one special character")
	}

	return nil
}