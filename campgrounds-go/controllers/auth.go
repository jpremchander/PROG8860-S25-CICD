package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"campgrounds-go/config"
	"campgrounds-go/models"
)

func ShowRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/register.html", gin.H{
		"title": "Register",
	})
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.HTML(http.StatusBadRequest, "auth/register.html", gin.H{
			"error": err.Error(),
			"title": "Register",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "auth/register.html", gin.H{
			"error": "Failed to process password",
			"title": "Register",
		})
		return
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user already exists
	var existingUser models.User
	err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		c.HTML(http.StatusBadRequest, "auth/register.html", gin.H{
			"error": "User with this email already exists",
			"title": "Register",
		})
		return
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "auth/register.html", gin.H{
			"error": "Failed to create user",
			"title": "Register",
		})
		return
	}

	// Log user in
	session := sessions.Default(c)
	session.Set("userID", result.InsertedID.(primitive.ObjectID).Hex())
	session.Set("username", user.Username)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/campgrounds")
}

func ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/login.html", gin.H{
		"title": "Login",
	})
}

func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "auth/login.html", gin.H{
			"error": "Invalid email or password",
			"title": "Login",
		})
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.HTML(http.StatusUnauthorized, "auth/login.html", gin.H{
			"error": "Invalid email or password",
			"title": "Login",
		})
		return
	}

	// Log user in
	session := sessions.Default(c)
	session.Set("userID", user.ID.Hex())
	session.Set("username", user.Username)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/campgrounds")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusSeeOther, "/")
}
