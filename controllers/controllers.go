package controllers

import (
	"log"
	"net/http"
	"ringer/ai"
	"ringer/models"

	"github.com/gin-gonic/gin"
)

func Index(context *gin.Context) {

	var message models.Message

	if err := context.BindJSON(&message); err != nil {
		// Handle error if the JSON is malformed
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ai.Gpt(message.Message)

	if err != nil {
		log.Fatal("error occured while ai working")
	}

	context.JSON(http.StatusOK, gin.H{
		"message": response,
	})
}
