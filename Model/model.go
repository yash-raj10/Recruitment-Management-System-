package model

import (
	"context"
	"fmt"
	"log"
	"os"

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


var Collection1 *mongo.Collection

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
	fmt.Println("Instance is Ready")
}