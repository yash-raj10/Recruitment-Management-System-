package main

import (
	controller "RMS/Controller"
	model "RMS/Model"
	"RMS/middlware"

	"github.com/gin-gonic/gin"
)

func main(){

// Initialize the database
	model.InitDB()

// Gin Setup
	r := gin.Default()

	r.POST("/signup", controller.Signup)
	r.POST("/login", controller.Login)
	r.POST("/uploadResume", middlware.AuthMiddleware(), controller.UploadResume)
	r.POST("/admin/job", middlware.AuthMiddleware(), middlware.AuthMiddleware(), controller.CreateJob)
	r.GET("/admin/applicants", middlware.AuthMiddleware(), middlware.AuthMiddleware(), controller.GetApplicants)
	r.GET("/admin/applicant/:applicant_id",middlware.AuthMiddleware(), middlware.AuthMiddleware(), controller.GetApplicant)

	r.Run(":8080")
}