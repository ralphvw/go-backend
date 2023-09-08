package main

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	"ralph.com/gotest/db"
)

func main() {
  db.InitDb()

  router := gin.Default()

  fmt.Println("Router running on :8080")
  log.Fatal(router.Run(":8080"))
}
