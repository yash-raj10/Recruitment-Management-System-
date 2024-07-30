package model

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Name            string             `bson:"name"`
	Email           string             `bson:"email"`
	Address         string             `bson:"address"`
	UserType        string             `bson:"userType"`
	PasswordHash    string             `bson:"passwordHash"`
	ProfileHeadline string             `bson:"profileHeadline"`
	Profile         Profile            `bson:"profile,omitempty"`
}

type Profile struct {
	ResumeFileAddress string `bson:"resumeFileAddress"`
	Skills            string `bson:"skills"`
	Education         string `bson:"education"`
	Experience        string `bson:"experience"`
	Name              string `bson:"name"`
	Email             string `bson:"email"`
	Phone             string `bson:"phone"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	UserID   string `json:"userId"`
	UserType string `json:"userType"`
	jwt.StandardClaims
}

type Job struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Title             string             `bson:"title"`
	Description       string             `bson:"description"`
	PostedOn          time.Time          `bson:"postedOn"`
	TotalApplications int                `bson:"totalApplications"`
	CompanyName       string             `bson:"companyName"`
	PostedBy          primitive.ObjectID `bson:"postedBy"`
}


var Collection1 *mongo.Collection
var Collection2 *mongo.Collection


func InitDB() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	Link := os.Getenv("Link")

	clientOption := options.Client().ApplyURI(Link)

	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection Success")

	Collection1 = client.Database("recruitment_system").Collection("users")
	Collection2 = client.Database("recruitment_system").Collection("jobs")

	fmt.Println("Instance is Ready")
}