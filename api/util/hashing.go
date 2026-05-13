package util

import (
    bcrypt "golang.org/x/crypto/bcrypt"
    
    "crypto/sha256"
    "encoding/hex"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func CreateTripcode(input string, strong bool) string {
    if strong {
        bytes, _ := bcrypt.GenerateFromPassword([]byte(input), 14)
        shortHash := bytes[:10]
        return hex.EncodeToString(shortHash)
    } else {
        fullHash := sha256.Sum256([]byte(input + "h-/-v9e/8==h-=f298h"))
        shortHash := fullHash[:10]
        return hex.EncodeToString(shortHash)
    }
}