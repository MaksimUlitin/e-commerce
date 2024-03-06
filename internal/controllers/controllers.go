package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/maksimulitin/internal/db"
	"github.com/maksimulitin/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	Validate                            = validator.New(validator.WithRequiredStructEnabled())
	UserCollection    *mongo.Collection = db.AbuotUser(db.Client, "Users")
	ProductCollection *mongo.Collection = db.AbuotProduct(db.Client, "Products")
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(enterPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(enterPassword), []byte(givenPassword))
	verify := true
	message := ""
	if err != nil {
		message = "Login or Password invalid"
		verify = false
	}
	return verify, message
}

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := Validate.Struct(user)

		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		count, err := UserCollection.CountDocuments(c, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
			return
		}

		count, err = UserCollection.CountDocuments(c, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "phone is already in use"})
			return
		}

	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}

}

func AddProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}

}

func ViewProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}

}

func SerchProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}

}
