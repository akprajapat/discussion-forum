package handlers

import (
	"time"

	"github.com/akprajapat/tic-tac-toe/discussion-forum/backend/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var db *mongo.Database

func Init(database *mongo.Database) {
	db = database
}

// --- AUTH ---

func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
	}
	_, err := db.Collection("users").InsertOne(c, user)
	if err != nil {
		c.JSON(500, gin.H{"error": "User exists"})
		return
	}
	c.JSON(200, gin.H{"message": "Registered"})
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	var user models.User
	err := db.Collection("users").FindOne(c, bson.M{"email": req.Email}).Decode(&user)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("secret"))
	c.JSON(200, gin.H{"token": tokenString})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if len(auth) < 8 {
			c.AbortWithStatusJSON(401, gin.H{"error": "No token"})
			return
		}
		tokenStr := auth[7:]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["user_id"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}
	}
}

// --- QUESTIONS ---

func GetQuestions(c *gin.Context) {
	q := c.Query("q")
	filter := bson.M{}
	if q != "" {
		filter = bson.M{"title": bson.M{"$regex": q, "$options": "i"}}
	}
	cur, _ := db.Collection("questions").Find(c, filter)
	var questions []models.Question
	cur.All(c, &questions)
	c.JSON(200, questions)
}

func CreateQuestion(c *gin.Context) {
	var req struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	userID, _ := primitive.ObjectIDFromHex(c.GetString("user_id"))
	question := models.Question{
		Title:     req.Title,
		Body:      req.Body,
		Author:    userID,
		Votes:     0,
		CreatedAt: time.Now(),
		Answers:   []primitive.ObjectID{},
	}
	res, _ := db.Collection("questions").InsertOne(c, question)
	question.ID = res.InsertedID.(primitive.ObjectID)
	c.JSON(200, question)
}

func GetQuestion(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var question models.Question
	err := db.Collection("questions").FindOne(c, bson.M{"_id": id}).Decode(&question)
	if err != nil {
		c.JSON(404, gin.H{"error": "Not found"})
		return
	}
	// Populate answers
	cur, _ := db.Collection("answers").Find(c, bson.M{"question": id})
	var answers []models.Answer
	cur.All(c, &answers)
	c.JSON(200, gin.H{"question": question, "answers": answers})
}

func VoteQuestion(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var req struct {
		Up bool `json:"up"`
	}
	c.ShouldBindJSON(&req)
	delta := 1
	if !req.Up {
		delta = -1
	}
	db.Collection("questions").UpdateOne(c, bson.M{"_id": id}, bson.M{"$inc": bson.M{"votes": delta}})
	c.JSON(200, gin.H{"message": "Voted"})
}

// --- ANSWERS ---

func CreateAnswer(c *gin.Context) {
	qid, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var req struct {
		Body string `json:"body"`
	}
	c.ShouldBindJSON(&req)
	userID, _ := primitive.ObjectIDFromHex(c.GetString("user_id"))
	answer := models.Answer{
		Body:      req.Body,
		Author:    userID,
		Question:  qid,
		Votes:     0,
		Comments:  []primitive.ObjectID{},
		CreatedAt: time.Now(),
	}
	res, _ := db.Collection("answers").InsertOne(c, answer)
	db.Collection("questions").UpdateOne(c, bson.M{"_id": qid}, bson.M{"$push": bson.M{"answers": res.InsertedID}})
	c.JSON(200, answer)
}

func VoteAnswer(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var req struct {
		Up bool `json:"up"`
	}
	c.ShouldBindJSON(&req)
	delta := 1
	if !req.Up {
		delta = -1
	}
	db.Collection("answers").UpdateOne(c, bson.M{"_id": id}, bson.M{"$inc": bson.M{"votes": delta}})
	c.JSON(200, gin.H{"message": "Voted"})
}

// --- COMMENTS ---

func CreateComment(c *gin.Context) {
	aid, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var req struct {
		Body string `json:"body"`
	}
	c.ShouldBindJSON(&req)
	userID, _ := primitive.ObjectIDFromHex(c.GetString("user_id"))
	comment := models.Comment{
		Body:      req.Body,
		Author:    userID,
		Answer:    aid,
		CreatedAt: time.Now(),
	}
	res, _ := db.Collection("comments").InsertOne(c, comment)
	db.Collection("answers").UpdateOne(c, bson.M{"_id": aid}, bson.M{"$push": bson.M{"comments": res.InsertedID}})
	c.JSON(200, comment)
}
