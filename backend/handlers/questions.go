package handlers

import (
	"discussion-forum/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
