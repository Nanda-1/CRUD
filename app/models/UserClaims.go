package models

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
