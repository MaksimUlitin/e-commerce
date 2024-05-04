package jwt

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt" // "github.com/dgrijalva/jwt-go"
	"github.com/maksimulitin/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type SingupDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.StandardClaims
}

var (
	AbuotUser  *mongo.Collection = db.AbuotUser(db.Client, "user")
	SEKRET_KEY                   = os.Getenv("SEKRET_KEY")
)

func GenerateToken(email string, firstname string, lastname string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SingupDetails{
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	RefreshToken := &SingupDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SEKRET_KEY))
	if err != nil {
		log.Println(err)
		return
	}

	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshToken).SignedString([]byte(SEKRET_KEY))
	if err != nil {
		log.Println(err)
		return
	}
	return token, refreshtoken, err
}

func ValidateToken(signedtoken string) (claims *SingupDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedtoken, SingupDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SEKRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SingupDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
	}

	return claims, msg

}

func UpdateToken() {

}
