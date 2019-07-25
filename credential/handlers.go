package credential

import (
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

// Credentials defines the structure of a user in the database
type Credentials struct {
	Username string
	Password string
}

// Claims defines the structure of claim type
type Claims struct {
	Username string
	jwt.StandardClaims
}