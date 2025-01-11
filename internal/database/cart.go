package database

import (
	"context"
	"errors"
	"github.com/maksimulitin/internal/models"
	"github.com/maksimulitin/lib/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"time"
)

var (
	ErrCantFindProduct    = errors.New("can't find product")
	ErrCantDecodeProducts = errors.New("can't decode products")
	ErrUserIDIsNotValid   = errors.New("user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add product to cart")
	ErrCantRemoveItem     = errors.New("cannot remove item from cart")
	ErrCantGetItem        = errors.New("cannot get item from cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	logger.Info("adding product to cart", slog.Any("productID", productID), slog.String("userID", userID))

	searchFromDb, err := prodCollection.Find(ctx, bson.M{"_id": productID})

	if err != nil {
		logger.Error("error finding product", slog.Any("productID", productID))
		return ErrCantFindProduct
	}

	var productCart []models.ProductUser
	err = searchFromDb.All(ctx, &productCart)

	if err != nil {
		logger.Error("error decoding products", slog.Any("productID", productID))
		return ErrCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		logger.Error("invalid user ID", slog.String("userID", userID))
		return ErrUserIDIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		logger.Error("error updating user cart", slog.Any("productID", productID), slog.String("userID", userID))
		return ErrCantUpdateUser
	}

	logger.Info("product added to cart successfully", slog.Any("productID", productID), slog.String("userID", userID))
	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	logger.Info("Removing item from cart", slog.Any("productID", productID), slog.String("userID", userID))

	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		logger.Error("Invalid user ID", slog.String("userID", userID))
		return ErrUserIDIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)

	if err != nil {
		logger.Error("Error removing item from cart", slog.Any("productID", productID), slog.String("userID", userID))
		return ErrCantRemoveItem
	}

	logger.Info("Item removed from cart successfully", slog.Any("productID", productID), slog.String("userID", userID))
	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	logger.Info("buying items from cart", slog.String("userID", userID))
	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		logger.Error("invalid user ID", slog.String("userID", userID))
		return ErrUserIDIsNotValid
	}

	var (
		getCartItems models.User
		orderCart    models.Order
	)

	orderCart.OrderID = primitive.NewObjectID()
	orderCart.OrderedAt = time.Now()
	orderCart.OrderCart = make([]models.ProductUser, 0)
	orderCart.PaymentMethod.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

	currentResults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})

	if err != nil {
		logger.Error("error aggregating cart items", slog.String("userID", userID))
		return ErrCantBuyCartItem
	}

	var getUserCart []bson.M

	if err = currentResults.All(ctx, &getUserCart); err != nil {
		logger.Error("error decoding aggregated results", slog.String("userID", userID))
		return ErrCantBuyCartItem
	}

	var totalPrice int32

	for _, userItem := range getUserCart {
		price := userItem["total"]
		totalPrice = price.(int32)
	}
	orderCart.Price = int(totalPrice)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)

	if err != nil {
		logger.Error("error updating user orders", slog.String("userID", userID))
		return ErrCantBuyCartItem
	}

	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems)

	if err != nil {
		logger.Error("error fetching user cart items", slog.String("userID", userID))
		return ErrCantBuyCartItem
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getCartItems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)

	if err != nil {
		logger.Error("error updating order list", slog.String("userID", userID))
		return ErrCantBuyCartItem
	}

	userCartEmpty := make([]models.ProductUser, 0)
	filtered := bson.D{primitive.E{Key: "_id", Value: id}}
	updated := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: userCartEmpty}}}}
	_, err = userCollection.UpdateOne(ctx, filtered, updated)

	if err != nil {
		logger.Error("error clearing user cart", slog.String("userID", userID))
		return ErrCantBuyCartItem
	}

	logger.Info("cart items purchased successfully", slog.String("userID", userID))
	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	logger.Info("instant buying product", slog.Any("productID", productID), slog.String("userID", userID))
	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		logger.Error("Invalid user ID", slog.String("userID", userID))
		return ErrUserIDIsNotValid
	}

	var (
		productDetails models.ProductUser
		ordersDetail   models.Order
	)
	ordersDetail.OrderID = primitive.NewObjectID()
	ordersDetail.OrderedAt = time.Now()
	ordersDetail.OrderCart = make([]models.ProductUser, 0)
	ordersDetail.PaymentMethod.COD = true

	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&productDetails)
	if err != nil {
		logger.Error("error fetching product details", slog.Any("productID", productID))
		return ErrCantFindProduct
	}

	ordersDetail.Price = productDetails.Price
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordersDetail}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		logger.Error("error updating user orders", slog.Any("productID", productID), slog.String("userID", userID))
		return ErrCantUpdateUser
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": productDetails}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)

	if err != nil {
		logger.Error("error updating order list", slog.Any("productID", productID), slog.String("userID", userID))
		return ErrCantUpdateUser
	}

	logger.Info("product purchased instantly", slog.Any("productID", productID), slog.String("userID", userID))
	return nil
}
