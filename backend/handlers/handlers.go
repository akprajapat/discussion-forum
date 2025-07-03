package handlers

import (
	"log"
	"time"

	"discussion-forum/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var db *mongo.Database

func Init(database *mongo.Database) {
	db = database
}

// --- AUTH ---

func Register(c *gin.Context) {
	log.Println("Register: called")
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Register: Invalid input:", err)
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Check for existing username or email
	var existing models.User
	err := db.Collection("users").FindOne(c, bson.M{
		"$or": []bson.M{
			{"username": req.Username},
			{"email": req.Email},
		},
	}).Decode(&existing)
	if err == nil {
		log.Println("Register: Username or email already exists")
		c.JSON(400, gin.H{"error": "Username or email already exists"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Register: Password hash error:", err)
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
	}
	_, err = db.Collection("users").InsertOne(c, user)
	if err != nil {
		log.Println("Register: DB insert error:", err)
		c.JSON(500, gin.H{"error": "User exists or DB error"})
		return
	}
	log.Println("Register: User registered:", req.Email)
	c.JSON(200, gin.H{"message": "Registered"})
}

func Login(c *gin.Context) {
	log.Println("Login: called")
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Login: Invalid input:", err)
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	var user models.User
	err := db.Collection("users").FindOne(c, bson.M{"username": req.Username}).Decode(&user)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		log.Println("Login: Invalid credentials for", req.Username)
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("secret"))
	log.Println("Login: Success for", req.Username)
	c.JSON(200, gin.H{"token": tokenString})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("AuthMiddleware: called")
		auth := c.GetHeader("Authorization")
		if len(auth) < 8 {
			log.Println("AuthMiddleware: No token")
			c.AbortWithStatusJSON(401, gin.H{"error": "No token"})
			return
		}
		tokenStr := auth[7:]
		token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Println("AuthMiddleware: Token valid for user_id", claims["user_id"])
			c.Set("user_id", claims["user_id"])
			c.Next()
		} else {
			log.Println("AuthMiddleware: Invalid token")
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}
	}
}
