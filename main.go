package main

import (
	controller "RMS/Controller"
	model "RMS/Model"

	"github.com/gin-gonic/gin"
)

func main(){

// Initialize the database
	model.InitDB()

// Gin Setup
	r := gin.Default()

	r.POST("/signup", controller.Signup)
	r.POST("/login", controller.Login)


	r.Run(":8080")
}