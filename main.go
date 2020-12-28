package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ginfolderstructure/api"
	"github.com/ginfolderstructure/database"
)

func main() {
	r := gin.New()
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/ping"},
	}))
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "DELETE", "HEAD"},
	}))
	r.Use(gin.Recovery())

	v1 := r.Group("/v1")
	v1.GET("/tiles/:service/:uuid/:timestamp/:z/:x/:y", api.GetTiles)
	v1.GET("/summary/:service/:uuid/:timestamp", api.GetSummary)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	database.Init()
	r.Run(":8081")
	database.Close()
}