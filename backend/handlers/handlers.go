package handlers

import (
	"log"
	"time"

	"discussion-forum/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// --- QUESTIONS ---

func GetQuestions(c *gin.Context) {
	log.Println("GetQuestions: called")
	q := c.Query("q")
	filter := bson.M{}
	if q != "" {
		filter = bson.M{"title": bson.M{"$regex": q, "$options": "i"}}
	}
	cur, err := db.Collection("questions").Find(c, filter)
	if err != nil {
		log.Println("GetQuestions: DB error:", err)
		c.JSON(500, gin.H{"error": "DB error"})
		return
	}
	var questions []models.Question
	if err := cur.All(c, &questions); err != nil {
		log.Println("GetQuestions: Cursor error:", err)
		c.JSON(500, gin.H{"error": "DB error"})
		return
	}
	log.Println("GetQuestions: returned", len(questions), "questions")
	c.JSON(200, questions)
}

func CreateQuestion(c *gin.Context) {
	log.Println("CreateQuestion: called")
	var req struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("CreateQuestion: Invalid input:", err)
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
	res, err := db.Collection("questions").InsertOne(c, question)
	if err != nil {
		log.Println("CreateQuestion: DB insert error:", err)
		c.JSON(500, gin.H{"error": "DB error"})
		return
	}
	question.ID = res.InsertedID.(primitive.ObjectID)
	log.Println("CreateQuestion: Question created with ID", question.ID.Hex())
	c.JSON(200, question)
}

func GetQuestion(c *gin.Context) {
	log.Println("GetQuestion: called")
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var question models.Question
	err := db.Collection("questions").FindOne(c, bson.M{"_id": id}).Decode(&question)
	if err != nil {
		log.Println("GetQuestion: Not found:", err)
		c.JSON(404, gin.H{"error": "Not found"})
		return
	}
	// Populate answers
	cur, _ := db.Collection("answers").Find(c, bson.M{"question": id})
	var answers []models.Answer
	cur.All(c, &answers)
	log.Println("GetQuestion: Found question", id.Hex())
	c.JSON(200, gin.H{"question": question, "answers": answers})
}

func VoteQuestion(c *gin.Context) {
	log.Println("VoteQuestion: called")
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var req struct {
		Up bool `json:"up"`
	}
	c.ShouldBindJSON(&req)
	delta := 1
	if !req.Up {
		delta = -1
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("VoteQuestion: Invalid input:", err)
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	_, err := db.Collection("questions").UpdateOne(c, bson.M{"_id": id}, bson.M{"$inc": bson.M{"votes": delta}})
	if err != nil {
		log.Println("VoteQuestion: DB update error:", err)
		c.JSON(500, gin.H{"error": "DB error"})
		return
	}
	log.Println("VoteQuestion: Voted", delta, "for question", id.Hex())
	c.JSON(200, gin.H{"message": "Voted"})
}

// --- ANSWERS ---

func CreateAnswer(c *gin.Context) {
	log.Println("CreateAnswer: called")
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
	res, err := db.Collection("answers").InsertOne(c, answer)
	if err != nil {
		log.Println("CreateAnswer: DB insert error:", err)
		c.JSON(500, gin.H{"error": "DB error"})
		return
	}
	_, err = db.Collection("questions").UpdateOne(c, bson.M{"_id": qid}, bson.M{"$push": bson.M{"answers": res.InsertedID}})
	if err != nil {
		log.Println("CreateAnswer: DB update error:", err)
		c.JSON(500, gin.H{"error": "DB error"})
		return
	}
	log.Println("CreateAnswer: Answer created with ID", res.InsertedID)
	c.JSON(200, answer)
}

func VoteAnswer(c *gin.Context) {
	log.Println("VoteAnswer: called")
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var req struct {
		Up bool `json:"up"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("VoteAnswer: Invalid input:", err)
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	delta := 1
	if !req.Up {
		delta = -1
	}
	// Update and return the new vote count
	res := db.Collection("answers").FindOneAndUpdate(
		c,
		bson.M{"_id": id},
		bson.M{"$inc": bson.M{"votes": delta}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var updated models.Answer
	if err := res.Decode(&updated); err != nil {
		log.Println("VoteAnswer: DB update error:", err)
		c.JSON(500, gin.H{"error": "DB error"})
		return
	}
	log.Println("VoteAnswer: Voted", delta, "for answer", id.Hex(), "new votes:", updated.Votes)
	c.JSON(200, gin.H{"votes": updated.Votes})
}

// --- COMMENTS ---

func CreateComment(c *gin.Context) {
	log.Println("CreateComment: called")
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
	res, err := db.Collection("comments").InsertOne(c, comment)
	if err != nil {
		log.Println("CreateComment: DB insert error:", err)
		c.JSON(500, gin.H{"error": "DB error"})
		return
	}
	_, err = db.Collection("answers").UpdateOne(c, bson.M{"_id": aid}, bson.M{"$push": bson.M{"comments": res.InsertedID}})
	if err != nil {
		log.Println("CreateComment: DB update error:", err)
		c.JSON(500, gin.H{"error": "DB error"})
		return
	}
	log.Println("CreateComment: Comment created with ID", res.InsertedID)
	c.JSON(200, comment)
}

func GetComments(c *gin.Context) {
	log.Println("GetComments: called")
	aid, _ := primitive.ObjectIDFromHex(c.Param("id"))
	cur, err := db.Collection("comments").Find(c, bson.M{"answer": aid})
	if err != nil {
		log.Println("GetComments: DB error:", err)
		c.JSON(500, gin.H{"error": "DB error"})
		return
	}
	var comments []models.Comment
	if err := cur.All(c, &comments); err != nil {
		log.Println("GetComments: Cursor error:", err)
		c.JSON(500, gin.H{"error": "DB error"})
		return
	}
	log.Println("GetComments: returned", len(comments), "comments")
	c.JSON(200, comments)
}
