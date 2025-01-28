package utils
import "golang.org/x/crypto/bcrypt"
func hashPassword(password string) (string, error) {
  
	// Hashing the password with a cost factor of 10
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func verifyPassword(hashedPassword, password string) bool {
	// Compare the hashed password with the given password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// If there is an error, the password does not match
		return false
	}
	// If no error, the password matches
	return true
}
