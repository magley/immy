package util

import (
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/crypto/argon2"
    
    "crypto/sha256"
)

// TODO: Don't hardcode.
var SECURE_TRIP_SALT = []byte("FH(3hf09ho3hcIHWFKUn2-=fw-0g4")

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
        bytes := argon2.IDKey([]byte(input), SECURE_TRIP_SALT, 3, 64 * 1024, 4, uint32(16))
        shortHash := bytes[:16]
        return "!!" + randomStringFromKey(shortHash, TRIP_CHARSET, 11)
    } else {
        fullHash := sha256.Sum256([]byte(input + "h-/-v9e/8==h-=f298h"))
        shortHash := fullHash[:16]
        return "!" + randomStringFromKey(shortHash, TRIP_CHARSET, 11)
    }
}

func randomStringFromKey(key []byte, charset string, length int) string {
    result := make([]byte, length)

    for i := 0; i < length; i++ {
        idx := key[i % len(key)] % byte(len(charset))
        result[i] = charset[idx]
    }

    return string(result)
}

func CreateUserID(ip string, threadID uint) string {
    const USERID_CHARSET = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890-=+/"
    fullHash := sha256.Sum256([]byte(ip + string(threadID) + "vhe98wf2*(Hfh839f3h8 C_+Wf3hci32=v2"))
    shortHash := fullHash[:16]
    return randomStringFromKey(shortHash, USERID_CHARSET, 7)
}