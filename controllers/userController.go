package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go_ecommerce_cart/database"
	"go_ecommerce_cart/helpers"
	"go_ecommerce_cart/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = helpers.OpenCollection(database.Client, "Users")
var ProductCollection *mongo.Collection = helpers.OpenCollection(database.Client, "Products")
var validate = validator.New()

func HashPassword(password string) string {
	byte, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(byte)
}

func VerifyPassword(userpassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Wrong credentials"
		valid = false
	}
	return valid, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validatorErr := validate.Struct(user)
		if validatorErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validatorErr.Error()})
			return
		}
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already in use"})
			return
		}
		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Phone already in use"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		token, refresh_token, _ := helpers.TokenGenerate(*user.Email, *user.First_Name, *user.Last_Name, *&user.User_ID, *user.User_Type)
		user.Token = &token
		user.Refresh_token = &refresh_token
		user.User_Cart = make([]models.Product, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, "Successflly signed up")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or email incorrect"})
			return
		}
		IsPasswordValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !IsPasswordValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refreshToken, _ := helpers.TokenGenerate(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, foundUser.User_ID, *foundUser.User_Type)
		defer cancel()
		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_ID)
		c.JSON(http.StatusFound, foundUser)
	}
}

//Search Product

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		var productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Something Went Wrong please after sometime")
			return
		}
		err = cursor.All(ctx, &productlist)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, productlist)
	}
}

func AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Product
		defer cancel()

		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		products.Product_ID = primitive.NewObjectID()
		products.User_ID = c.GetString("uid")
		_, anyerr := ProductCollection.InsertOne(ctx, products)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Successfully added the product")
	}
}

func SearchByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProducts []models.Product
		queryParam := c.Query("name")
		if queryParam == "" {
			log.Println("no Product found")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invaliad Search Index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchQueryDB, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(404, "Failed to fetch the query")
			return
		}
		err = searchQueryDB.All(ctx, &searchProducts)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchQueryDB.Close(ctx)
		if err := searchQueryDB.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchProducts)
	}
}
