package jwt

import "github.com/golang-jwt/jwt" // "github.com/dgrijalva/jwt-go"

type FFDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.StandardClaims
}

var ()

func GenerateToken() {

}

func ValidateToken() {

}

func UpdateToken() {

}
