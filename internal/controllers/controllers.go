package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/maksimulitin/internal/database"
	"github.com/maksimulitin/internal/models"
	generate "github.com/maksimulitin/internal/tokens"
	"github.com/maksimulitin/lib/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"time"
)

var (
	UserCollection    *mongo.Collection = database.UserData(database.Client, "Users")
	ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
	Validate                            = validator.New()
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		logger.Error("Error hashing password", slog.Any("error", err))
		panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, thisPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(thisPassword), []byte(userPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "Login Or Password is Incorrect"
		valid = false
		logger.Error("Password verification failed", slog.Any("error", err))
	}

	return valid, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			logger.Error("Error binding JSON", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := Validate.Struct(user)

		if validationErr != nil {
			logger.Error("Validation failed", slog.Any("error", validationErr))
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			logger.Error("Error counting documents", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			logger.Info("User already exists", slog.String("email", *user.Email))
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		if err != nil {
			logger.Error("Error counting documents by phone", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			logger.Info("Phone is already in use", slog.String("phone", *user.Phone))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone is already in use"})
			return
		}

		hashedPassword := HashPassword(*user.Password)
		user.Password = &hashedPassword

		token, refreshToken, _ := generate.TokenGenerator(*user.Email, *user.FirstName, *user.LastName, user.UserID)
		user.Token = &token
		user.RefreshToken = &refreshToken

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()

		user.UserCart = make([]models.ProductUser, 0)
		user.AddressDetails = make([]models.Address, 0)
		user.OrderStatus = make([]models.Order, 0)

		_, insertErr := UserCollection.InsertOne(ctx, user)

		if insertErr != nil {
			logger.Error("Error inserting user", slog.Any("error", insertErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}

		logger.Info("User successfully signed up", slog.String("userID", user.UserID))
		c.JSON(http.StatusCreated, "Successfully Signed Up!!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var (
			user      models.User
			foundUser models.User
		)

		if err := c.BindJSON(&user); err != nil {
			logger.Error("Error binding JSON", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if err != nil {
			logger.Error("Error finding user", slog.Any("email", user.Email), slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}

		PasswordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		if !PasswordIsValid {
			logger.Error("Invalid password", slog.Any("email", user.Email))
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := generate.TokenGenerator(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserID)
		generate.UpdateAllTokens(token, refreshToken, foundUser.UserID)

		logger.Info("User logged in successfully", slog.String("userID", foundUser.UserID))
		c.JSON(http.StatusFound, foundUser)
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var products models.Product

		if err := c.BindJSON(&products); err != nil {
			logger.Error("Error binding JSON", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		products.ProductID = primitive.NewObjectID()
		_, anyErr := ProductCollection.InsertOne(ctx, products)

		if anyErr != nil {
			logger.Error("Error inserting product", slog.Any("error", anyErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}

		logger.Info("Product successfully added", slog.String("productID", products.ProductID.Hex()))
		c.JSON(http.StatusOK, "Successfully added our Product Admin!!")
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			productList []models.Product
			ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})

		if err != nil {
			logger.Error("Error finding products", slog.Any("error", err))
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong, please try again later")
			return
		}

		err = cursor.All(ctx, &productList)

		if err != nil {
			logger.Error("Error reading cursor", slog.Any("error", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil {
			logger.Error("Cursor error", slog.Any("error", err))
			c.IndentedJSON(http.StatusBadRequest, "Invalid")
			return
		}

		logger.Info("Products fetched successfully", slog.Int("count", len(productList)))
		c.IndentedJSON(http.StatusOK, productList)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProducts []models.Product
		queryParam := c.Query("name")

		if queryParam == "" {
			logger.Warn("Empty query parameter")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchDbQuery, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})

		if err != nil {
			logger.Error("Error querying database", slog.Any("query", queryParam), slog.Any("error", err))
			c.IndentedJSON(http.StatusNotFound, "Something went wrong in fetching the database query")
			return
		}

		err = searchDbQuery.All(ctx, &searchProducts)

		if err != nil {
			logger.Error("Error reading query results", slog.Any("error", err))
			c.IndentedJSON(http.StatusBadRequest, "Invalid")
			return
		}

		defer searchDbQuery.Close(ctx)

		if err := searchDbQuery.Err(); err != nil {
			logger.Error("Query cursor error", slog.Any("error", err))
			c.IndentedJSON(http.StatusBadRequest, "Invalid request")
			return
		}

		logger.Info("Products fetched successfully by query", slog.String("query", queryParam), slog.Int("count", len(searchProducts)))
		c.IndentedJSON(http.StatusOK, searchProducts)
	}
}
