package util

import (
	"os"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"

	"crypto/sha256"
)

var SECURE_TRIP_SALT = []byte(os.Getenv("SECURE_TRIP_SALT"))
var INSECURE_TRIP_SALT = os.Getenv("INSECURE_TRIP_SALT")
var USERID_SALT = os.Getenv("USERID_SALT")

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateTripcode(input string, strong bool) string {
	const TRIP_CHARSET = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890-=+"

	if strong {
		bytes := argon2.IDKey([]byte(input), SECURE_TRIP_SALT, 3, 64*1024, 4, uint32(16))
		shortHash := bytes[:16]
		return "!!" + randomStringFromKey(shortHash, TRIP_CHARSET, 11)
	} else {
		fullHash := sha256.Sum256([]byte(input + INSECURE_TRIP_SALT))
		shortHash := fullHash[:16]
		return "!" + randomStringFromKey(shortHash, TRIP_CHARSET, 11)
	}
}

func randomStringFromKey(key []byte, charset string, length int) string {
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		idx := key[i%len(key)] % byte(len(charset))
		result[i] = charset[idx]
	}

	return string(result)
}

func CreateUserID(ip string, threadID uint) string {
	const USERID_CHARSET = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890-=+/"
	fullHash := sha256.Sum256([]byte(ip + string(threadID) + USERID_SALT))
	shortHash := fullHash[:16]
	return randomStringFromKey(shortHash, USERID_CHARSET, 7)
}
