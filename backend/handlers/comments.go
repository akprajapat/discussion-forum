package handlers

import (
	"discussion-forum/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
