package controllers

import "golang.org/x/crypto/bcrypt"

// HashPassword takes a plain text password and returns a bcrypt hashed version
func HashPassword(password string) (string, error) {
	// Generate a salt with default cost (10)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword compares a plain text password with a hashed password
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
