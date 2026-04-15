package routes

import (
	controller "go_ecommerce_cart/controllers"
	"go_ecommerce_cart/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/user/signup", controller.SignUp())
	incomingRoutes.POST("/user/login", controller.Login())

	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/user/productview", controller.SearchProduct())
	incomingRoutes.POST("/admin/addproduct", controller.AddProduct())
	incomingRoutes.GET("/user/search", controller.SearchByQuery())
}
