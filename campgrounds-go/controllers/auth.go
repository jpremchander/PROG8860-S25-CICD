package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"yelpcamp-go/config"
	"yelpcamp-go/models"
	"yelpcamp-go/utils"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// API Registration
func (ac *AuthController) Register(c *gin.Context) {
	var input models.UserInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()
	
	// Check if user already exists
	if _, err := models.FindUserByUsername(db, input.Username); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}
	
	if _, err := models.FindUserByEmail(db, input.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := user.Create(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	token, err := utils.GenerateJWTToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
		"user":    gin.H{"id": user.ID, "username": user.Username, "email": user.Email},
	})
}

// Web Registration - Integrated with auto JWT
func (ac *AuthController) RegisterWeb(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	log.Printf("🔐 Registration attempt - Username: %s, Email: %s", username, email)

	if username == "" || email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register - YelpCamp",
			"error": "All fields are required",
		})
		return
	}

	db := config.GetDB()
	
	// Check if user already exists
	if _, err := models.FindUserByUsername(db, username); err == nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register - YelpCamp",
			"error": "Username already exists",
		})
		return
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	if err := user.Create(db); err != nil {
		log.Printf("❌ Error creating user: %v", err)
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"title": "Register - YelpCamp",
			"error": "Could not create user",
		})
		return
	}

	// Use the integrated JWT generation
	token, err := utils.GenerateJWTToken(user.ID)
	if err != nil {
		log.Printf("❌ Error generating token: %v", err)
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"title": "Register - YelpCamp",
			"error": "Could not generate token",
		})
		return
	}

	log.Printf("✅ User registered successfully: %s", username)

	// Set secure cookie with proper settings
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", token, 3600*24*7, "/", "", false, true) // 7 days, httpOnly
	
	// Redirect to campgrounds page
	c.Redirect(http.StatusSeeOther, "/campgrounds")
}

// API Login
func (ac *AuthController) Login(c *gin.Context) {
	var input models.LoginInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDB()
	user, err := models.FindUserByUsername(db, input.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := user.CheckPassword(input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Use the integrated JWT generation
	token, err := utils.GenerateJWTToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user":    gin.H{"id": user.ID, "username": user.Username, "email": user.Email},
	})
}

// Web Login - Integrated with auto JWT
func (ac *AuthController) LoginWeb(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	log.Printf("🔐 Login attempt - Username: %s", username)

	if username == "" || password == "" {
		log.Printf("❌ Login failed - Missing credentials")
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"title": "Login - YelpCamp",
			"error": "Username and password are required",
		})
		return
	}

	db := config.GetDB()
	user, err := models.FindUserByUsername(db, username)
	if err != nil {
		log.Printf("❌ Login failed - User not found: %s", username)
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"title": "Login - YelpCamp",
			"error": "Invalid username or password",
		})
		return
	}

	if err := user.CheckPassword(password); err != nil {
		log.Printf("❌ Login failed - Invalid password for user: %s", username)
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"title": "Login - YelpCamp",
			"error": "Invalid username or password",
		})
		return
	}

	// Use the integrated JWT generation
	token, err := utils.GenerateJWTToken(user.ID)
	if err != nil {
		log.Printf("❌ Login failed - Token generation error: %v", err)
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"title": "Login - YelpCamp",
			"error": "Could not generate token",
		})
		return
	}

	log.Printf("✅ Login successful for user: %s (ID: %s)", username, user.ID.Hex())

	// Set secure cookie with proper settings
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", token, 3600*24*7, "/", "", false, true) // 7 days, httpOnly
	
	log.Printf("🍪 Cookie set for user: %s", username)
	
	// Redirect to campgrounds page
	c.Redirect(http.StatusSeeOther, "/campgrounds")
}

// Logout
func (ac *AuthController) Logout(c *gin.Context) {
	log.Printf("🚪 User logging out")
	// Clear the cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}
