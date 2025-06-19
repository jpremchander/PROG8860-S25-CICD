package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"campgrounds-app/config"
	"campgrounds-app/models"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

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

	token, err := generateToken(user.ID)
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

	token, err := generateToken(user.ID)
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

func generateToken(userID primitive.ObjectID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.Hex(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
