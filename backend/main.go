package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	controller "github.com/wardhana-org/ekstra/backend/controllers"
)

func main() {
	router := gin.Default()

	router.GET("/ping", controller.Ping)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}

}
