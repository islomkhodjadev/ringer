package routes

import (
	"ringer/controllers"

	"github.com/gin-gonic/gin"
)

func Routers(router *gin.Engine) {
	router.SetTrustedProxies([]string{"127.0.0.1"}) // Example for local development

	router.POST("/", controllers.Index)

}
