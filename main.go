package main

import (
	"ringer/database"
	"ringer/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connection()
	var router *gin.Engine = gin.Default()
	routes.Routers(router)
	router.Run(":8000")
}
