package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/maksimulitin/internal/db"
	"github.com/maksimulitin/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		var (
			user      models.User
			findUser  models.User
			c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		)
		defer cancel()

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := UserCollection.FindOne(c, bson.M{"email": user.Email}).Decode(&findUser)
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "login or password error"})
			return
		}

		passwordValet, msg := VerifyPassword(*user.Password, *findUser.Password)
		if !passwordValet {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

	}

}

func AdminProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			product   models.Porduct
			c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		)
		defer cancel()

		if err := ctx.BindJSON(&product); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product.ProductID = primitive.NewObjectID()
		_, anyErr := ProductCollection.InsertOne(c, &product)

		if anyErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}
		defer cancel()

		ctx.JSON(http.StatusOK, "successfully added our product admin!")
	}

}

func ViewProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			productList models.Porduct
			c, cancel   = context.WithTimeout(context.Background(), 100*time.Second)
		)
		defer cancel()
		cursor, err := ProductCollection.Find(c, bson.D{})

		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		err = cursor.All(c, &productList)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(c)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			ctx.IndentedJSON(http.StatusBadRequest, "invalid")
			return
		}
		defer cancel()
		ctx.IndentedJSON(200, productList)
	}

}

func SerchProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			productsSerch []models.Porduct
			c, cancel     = context.WithTimeout(context.Background(), 100*time.Second)
		)
		queryParam := ctx.Query("name")

		if queryParam == "" {
			log.Println("конч напиши чтото пжпж")
			ctx.Header("Content-Type", "Aplication/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "invalid serch index"})
		}
		defer cancel()
		serchBD, err := ProductCollection.Find(c, bson.M{"productName": bson.M{"$regex": queryParam}})
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, "something went wrong in fetching the dbquery")
			return
		}
		err = serchBD.All(c, &productsSerch)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer serchBD.Close(c)
		if err := serchBD.Err(); err != nil {
			log.Println(err)
			ctx.IndentedJSON(http.StatusBadRequest, "invalid request")
			return
		}
		defer cancel()
		ctx.IndentedJSON(http.StatusOK, serchBD)
	}

}
