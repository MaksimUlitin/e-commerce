package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName      *string            `json:"first_name" validate:"required,min=2,max=30"`
	LastName       *string            `json:"last_name" validate:"required,min=2,max=30"`
	Email          *string            `json:"email" validate:"email, required"`
	Phone          *string            `json:"phone" validate:"required"`
	Token          *string            `json:"token"`
	RefreshToken   *string            `json:"refresh_token"`
	Password       *string            `json:"password" validate:"required,min=8"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	UserID         *string            `json:"user_id"`
	AddressDetalis []Address          `json:"address_detalis" bson:"address_detalis"`
	OrderStatus    []Order            `json:"order_status" bson:"order_status"`
	UserCart       []ProductUser      `json:"user_cart" bson:"user_cart"`
}

type Porduct struct {
	ProductID   primitive.ObjectID `bson:"_id"`
	ProductName *string            `json:"Product_name"`
	Price       *uint64            `json:"price"`
	Rating      *uint8             `json:"rating"`
	Imges       *string            `json:"imges"`
}

type ProductUser struct {
	ProductID   primitive.ObjectID `bson:"_id"`
	ProductName *string            `json:"product_name" bson:"product_name"`
	Price       *uint32            `json:"price" bson:"price"`
	Rating      *uint8             `json:"rating" bson:"rating"`
	Imges       *string            `json:"imges" bson:"imges"`
}

type Address struct {
	AddressID primitive.ObjectID `bson:"_id"`
	Cite      *string            `json:"cite" bson:"cite"`
	Street    *string            `json:"street" bson:"srteet"`
	House     *string            `json:"house" bson:"house"`
	Pincode   *string            `json:"pincode" bson:"pincode"`
}

type Order struct {
	OrderID       primitive.ObjectID `bson:"_id"`
	OrderCart     []ProductUser      `json:"order_cart" bson:"order_cart"`
	OrdereredAt   time.Time          `json:"orderered_at" bson:"orderered_at"`
	Price         *uint64            `json:"price" bson:"price"`
	Discount      *int               `json:"discount" bson:"discount"`
	PaymentMethod Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"cod" bson:"cod"`
}
