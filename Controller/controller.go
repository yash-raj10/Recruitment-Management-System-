package controller

import (
	model "RMS/Model"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)


func Signup(c *gin.Context){
	var user model.User
	

	if err := c.BindJSON(&user)
	err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.PasswordHash = string(hashedPassword)

	_, err = model.Collection1.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User added successfully"})
}

func Login(c *gin.Context){
	if err := godotenv.Load(".env.local")
	err != nil {
		log.Fatalf("Error loading .env file")
	}
	Key := os.Getenv("JWT")
	var jwtKey = []byte(Key)

	var loginData model.LoginData

	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	err := model.Collection1.FindOne(context.TODO(), bson.M{"email": loginData.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found with email"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	expireTime := time.Now().Add(24 * time.Hour)
	claims := &model.Claims{
		UserID:   user.ID.Hex(),
		UserType: user.UserType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}


func UploadResume(c *gin.Context) {
	userType, _ := c.Get("userType")
	if userType != "Applicant" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only applicants can upload resumes"})
		return
	}

	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Check extension
	ext := filepath.Ext(file.Filename)
	if ext != ".pdf" && ext != ".docx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF & DOCX are allowed"})
		return
	}

	// Save file
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	if err := c.SaveUploadedFile(file, "uploads/"+filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	fileH, err := os.Open("uploads/" + filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to opennn resume"})
		return
	}
	defer fileH.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.apilayer.com/resume_parser/upload", fileH)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make api call (resume)"})
		return
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("apikey", "gNiXyflsFu3WNYCz1ZCxdWDb7oQg1Nl1")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get resp (resume)"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get body (resume)"})
		return
	}

	var resumeData ResumeData
	err = json.Unmarshal(body, &resumeData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get resumeData (resume)", "data": resp.Body})
		return
	}

	// Update user profile section with resume data
	userId, _ := c.Get("userId")
	objectId, _ := primitive.ObjectIDFromHex(userId.(string))
	_, err = model.Collection1.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectId},
		bson.M{"$set": bson.M{
			"profile.resumeFileAddress": "uploads/" + filename,
			"profile.skills":            strings.Join(resumeData.Skills, ", "),
			"profile.education":         resumeData.Education,
			"profile.experience":        resumeData.Experience,
			"profile.name":              resumeData.Name,
			"profile.email":             resumeData.Email,
			"profile.phone":             resumeData.Phone,
		}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resume uploaded and processed successfully"})
}

type ResumeData struct {
	Education  string `json:"education"`
	Email      string      `json:"email"`
	Experience string `json:"experience"`
	Name       string      `json:"name"`
	Phone      string      `json:"phone"`
	Skills     []string    `json:"skills"`
}

func CreateJob(c *gin.Context) {
	var job model.Job
	if err := c.BindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, _ := c.Get("userId")
	objectId, _ := primitive.ObjectIDFromHex(userId.(string))
	job.PostedBy = objectId
	job.PostedOn = time.Now()

	_, err := model.Collection2.InsertOne(context.TODO(), job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Job created successfully"})
}

func GetApplicants(c *gin.Context) {
	cursor, err := model.Collection1.Find(context.TODO(), bson.M{"userType": "Applicant"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applicants"})
		return
	}
	defer cursor.Close(context.TODO())

	var applicants []model.User
	if err = cursor.All(context.TODO(), &applicants); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode applicants"})
		return
	}

	c.JSON(http.StatusOK, applicants)
}

func GetApplicant(c *gin.Context) {
	applicantId, err := primitive.ObjectIDFromHex(c.Param("applicant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid applicant ID"})
		return
	}

	var applicant model.User
	err =  model.Collection1.FindOne(context.TODO(), bson.M{"_id": applicantId, "userType": "Applicant"}).Decode(&applicant)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Applicant not found"})
		return
	}

	c.JSON(http.StatusOK, applicant)
}