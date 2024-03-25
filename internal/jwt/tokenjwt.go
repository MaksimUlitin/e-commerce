package jwt

import (
	"os"

	"github.com/golang-jwt/jwt" // "github.com/dgrijalva/jwt-go"
	"github.com/maksimulitin/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type FFDetails struct {
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

func GenerateToken(email string, firstname string, lastname string, uid string) {

}

func ValidateToken() {

}

func UpdateToken() {

}
