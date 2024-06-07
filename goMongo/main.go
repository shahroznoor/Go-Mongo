package main

import (
	"example.com/goMongo/db"
	"example.com/goMongo/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":3000")
}
