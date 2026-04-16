package main

import (
	"go_ecommerce_cart/controllers"
	"go_ecommerce_cart/database"
	"go_ecommerce_cart/helpers"
	"go_ecommerce_cart/middleware"
	"go_ecommerce_cart/routes"

	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	app := controllers.NewApplication(helpers.OpenCollection(database.Client, "Products"),
		helpers.OpenCollection(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)

	router.Use(middleware.Authentication())
	router.GET("/cart/add", app.AddToCart())
	router.GET("/cart/remove", app.RemoveItem())
	router.GET("/cart", app.GetItemFromCart())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())
	router.POST("/address/add", controllers.AddAddress())
	router.PUT("/address/edithome", controllers.EditHomeAddress())
	router.PUT("/address/editwork", controllers.EditWorkAddress())
	router.GET("/address/delete", controllers.DeleteAddress())

	router.Run(":" + port)

}
