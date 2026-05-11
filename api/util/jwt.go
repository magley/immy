package util

import (
    "fmt"
    "time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("temporary-key") // TODO: Use private key pair, store keys in environment.

type JWTClaims struct {
    Username 	string 			`json:"username"`
    Id 			uint 			`json:"id"`
    Role 		string			`json:"role"`
    jwt.RegisteredClaims
}

func CreateJWT(userId uint, username string, userRole string) (string, error) {
    expiresAt := time.Now().Add(12 * time.Hour)

    claims := &JWTClaims{
        Id: userId,
        Username: username,
        Role: userRole,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateJWT(token string) (*JWTClaims, error) {
    var claims JWTClaims

    tokn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
        return jwtKey, nil
    })

    if err != nil {
        return nil, err
    }
    if !tokn.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return &claims, nil
}