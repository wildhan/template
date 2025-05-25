package helper

import "golang.org/x/crypto/bcrypt"

func EmptyStringToNull(input string) interface{} {
	if input != "" {
		return input
	} else {
		return nil
	}

}

// Fungsi untuk mengenkripsi password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Fungsi untuk membandingkan password dengan hash-nya
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
