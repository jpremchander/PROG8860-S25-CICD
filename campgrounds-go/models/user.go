package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"-" bson:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserInput struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) Create(db *mongo.Database) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	
	if err := u.HashPassword(); err != nil {
		return err
	}

	collection := db.Collection("users")
	result, err := collection.InsertOne(context.Background(), u)
	if err != nil {
		return err
	}
	
	u.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func FindUserByUsername(db *mongo.Database, username string) (*User, error) {
	var user User
	collection := db.Collection("users")
	
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func FindUserByEmail(db *mongo.Database, email string) (*User, error) {
	var user User
	collection := db.Collection("users")
	
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func FindUserByID(db *mongo.Database, id primitive.ObjectID) (*User, error) {
	var user User
	collection := db.Collection("users")
	
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}
