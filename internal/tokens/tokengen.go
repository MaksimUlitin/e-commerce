package token

import (
	"context"
	"github.com/maksimulitin/internal/database"
	"github.com/maksimulitin/lib/logger"
	"log/slog"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.StandardClaims
}

var (
	UserData   *mongo.Collection = database.UserData(database.Client, "Users")
	SECRET_KEY                   = os.Getenv("SECRET_LOVE")
)

func TokenGenerator(email string, firstname string, lastname string, uid string) (signedToken string, signedRefreshToken string, err error) {
	logger.Info("Generating tokens", slog.String("email", email), slog.String("uid", uid))

	claims := &SignedDetails{
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		logger.Error("Error generating token", slog.Any("error", err))
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		logger.Error("Error generating refresh token", slog.Any("error", err))
		return "", "", err
	}

	logger.Info("Tokens generated successfully", slog.String("email", email), slog.String("uid", uid))
	return token, refreshToken, nil
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	logger.Info("Validating token")

	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		logger.Error("Error parsing token", slog.Any("error", err))
		msg = err.Error()
		return nil, msg
	}

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		msg = "token is invalid"
		logger.Warn("Invalid token")
		return nil, msg
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		logger.Warn("Token expired")
		return nil, msg
	}

	logger.Info("Token validated successfully", slog.String("uid", claims.Uid))
	return claims, ""
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	logger.Info("Updating all tokens", slog.String("userId", userId))

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: updatedAt})
	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := UserData.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updateObj},
	}, &opt)

	if err != nil {
		logger.Error("Error updating tokens in database", slog.Any("error", err), slog.String("userId", userId))
		return
	}

	logger.Info("Tokens updated successfully", slog.String("userId", userId))
}
