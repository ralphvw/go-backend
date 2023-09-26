package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"ralph.com/gotest/db"
	"ralph.com/gotest/handlers"
	"ralph.com/gotest/services"
)

func main() {
	db := db.InitDb()

	router := gin.Default()

	router.POST("/users", func(c *gin.Context) {
		handlers.CreateUser(c, db)
	})

	router.GET("/users", func(c *gin.Context) {
		handlers.GetAllUsers(c, db)
	})

	router.POST("/auth/login", func(c *gin.Context) {
		services.LoginHandler(c, db)
	})

	fmt.Println("Router running on :8080")
	log.Fatal(router.Run(":8080"))
}
