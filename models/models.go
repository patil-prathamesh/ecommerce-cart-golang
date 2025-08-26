package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName      string             `json:"first_name" validate:"required,min=2,max=30"`
	LastName       string             `json:"last_name" validate:"required,min=2,max=30"`
	Password       string             `json:"password" validate:"required,min=6"`
	Email          string             `json:"email" validate:"required,email"`
	Phone          string             `json:"phone"  validate:"required"`
	RefreshToken   string             `json:"refresh_token" bson:"refresh_token"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	UserCart       []ProductUser      `json:"user_cart" bson:"user_cart"`
	AddressDetails []Address          `json:"address_details" bson:"address_details"`
	Order          []Order            `json:"orders" bson:"orders"`
}

type Product struct {
	ProductID   primitive.ObjectID `bson:"_id,omitempty"`
	ProductName string             `json:"product_name" bson:"product_name"`
	Price       uint64             `json:"price"`
	Rating      uint8              `json:"rating"`
	Image       string             `json:"image"`
}

type ProductUser struct {
	ProductID   primitive.ObjectID `bson:"_id,omitempty"`
	ProductName string             `json:"product_name" bson:"product_name"`
	Price       uint64             `json:"price" bson:"price"`
	Rating      uint8              `json:"rating" bson:"rating"`
	Image       string             `json:"image" bson:"image"`
}

type Address struct {
	AddressID primitive.ObjectID `bson:"_id"`
	House     string             `json:"house" bson:"house"`
	Street    string             `json:"street" bson:"street"`
	City      string             `json:"city" bson:"city"`
	Pincode   string             `json:"pin_code" bson:"pin_code"`
}

type Order struct {
	OrderID       primitive.ObjectID `bson:"_id"`
	OrderCart     []ProductUser      `json:"order_list" bson:"order_list"`
	OrderedAt     time.Time          `json:"order_at" bson:"order_at"`
	Price         uint64             `json:"price" bson:"price"`
	Discount      uint8              `json:"discount" bson:"discount"`
	PaymentMethod Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool
	COD     bool
}
