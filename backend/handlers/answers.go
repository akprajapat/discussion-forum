package handlers

import (
	"discussion-forum/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
