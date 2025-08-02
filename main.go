// main.go
package main

import (
	"net/http"
	"strings"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"github.com/joho/godotenv"
	"expense-tracker/db"
	"fmt"
	"expense-tracker/models"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET")) // replace with env var in production

type Expense struct {
	Amount   float64 `json:"amount"`
	Note     string  `json:"note"`
	Category string  `json:"category"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		// log.Fatal("❌ Error loading .env file")
	}
	db.InitDB()
	r := gin.Default()
	r.StaticFile("/", "./public/index.html")


	r.POST("/login", loginHandler)

	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.POST("/expense", expenseHandler)
	}

	r.Run(":8080")
}

func loginHandler(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Dummy validation (replace with DB check)
	if creds.Email != "test@gmail.com" || creds.Password != "test" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokenString, err := GenerateJWT(creds.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func expenseHandler(c *gin.Context) {
	var expense models.Expense

	if err := c.ShouldBindJSON(&expense); err != nil {
		fmt.Println("❌ Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🔍 Log the incoming request
	fmt.Println("✅ Expense received:", expense)

	// Save to DB
	if err := db.DB.Create(&expense).Error; err != nil {
		fmt.Println("❌ Failed to insert into DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save expense"})
		return
	}

	fmt.Println("✅ Expense inserted successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Expense added"})
}


func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user", claims.Subject)
		c.Next()
	}
}

func GenerateJWT(email string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}
