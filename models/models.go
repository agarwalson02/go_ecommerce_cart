package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id"`
	First_Name      *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_Name       *string            `json:"last_name" validate:"required,min=2,max=100"`
	Password        *string            `json:"password" validate:"required"`
	Email           *string            `json:"email" validate:"required"`
	Phone           *string            `json:"phone" validate:"required"`
	User_Type       *string            `json:"user_type" validate:"required"`
	Token           *string            `json:"token"`
	Refresh_token   *string            `json:"refresh_token"`
	Created_at      time.Time          `json:"createdat"`
	Updated_at      time.Time          `json:"updatedat"`
	User_ID         string             `json:"user_id"`
	User_Cart       []Product          `json:"usercart" bson:"user_cart"`
	Address_Details []Address          `json:"address" bson:"address"`
	Order_Status    []Order            `json:"orders" bson:"orders"`
}

type Product struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name"`
	Price        *float64           `json:"price"`
	Rating       *uint64            `json:"rating"`
	Image        *string            `json:"image"`
	User_ID      string             `json:"user_id"`
}

// type UserProduct struct {
// 	Product_ID   primitive.ObjectID `bson:"_id"`
// 	Product_Name *string            `json:"product_name" bson:"product_name"`
// 	Price        *float64           `json:"price" bson:"price"`
// 	Rating       *uint64            `json:"rating" bson:"rating"`
// 	Image        *string            `json:"image" bson:"image"`
// }

type Address struct {
	Address_ID primitive.ObjectID `bson:"_id"`
	House      *string            `json:"house_name" bson:"house_name"`
	Street     *string            `json:"street_name" bson:"street_name"`
	City       *string            `json:"city_name" bson:"city_name"`
	Pincode    *string            `json:"pin_code" bson:"pin_code"`
}

type Order struct {
	Order_Id       primitive.ObjectID `bson:"_id"`
	Order_Cart     []Product          `json:"order_list"  bson:"order_list"`
	Orderered_At   time.Time          `json:"ordered_on"  bson:"ordered_on"`
	Price          int                `json:"total_price" bson:"total_price"`
	Discount       *int               `json:"discount"    bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"COD" bson:"cod"`
}
