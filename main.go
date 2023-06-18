package main

import (
	"soccer-notifs/controllers"
	"soccer-notifs/initializers"
	"soccer-notifs/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.LoadDB()
	initializers.SyncDatabase()
}
func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/validate",middleware.IsAuthenticated,controllers.Validated)
	r.POST("/register",controllers.Register )
	r.POST("/login",controllers.Login )
	
	r.Run()
	
}
