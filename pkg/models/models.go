package models

import "time"

type User struct {
	ID             uint
	FirstName      *string
	LastName       *string
	Email          *string
	Phone          *string
	Token          *string
	RefreshToken   *string
	Password       *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	UserID         *string
	AddressDetalis []Address
	OrderStatus    []Order
	UserCart       []ProductUser
}

type Porduct struct {
	ProductID   uint
	ProductName *string
	Price       *uint64
	Rating      *uint8
	Imges       *string
}

type ProductUser struct {
	ProductID   uint
	ProductName *string
	Price       *uint32
	Rating      *uint8
	Imges       *string
}

type Address struct {
	AddressID uint
	Cite      *string
	Street    *string
	House     *string
	Pincode   *string
}

type Order struct {
	OrderID       uint
	OrderCart     []ProductUser
	OrdereredAt   time.Time
	Price         *uint64
	Discount      *int
	PaymentMethod Payment
}

type Payment struct {
	Digital bool
	COD     bool
}
