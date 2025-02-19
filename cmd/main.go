package main

import (
	"github.com/gin-gonic/gin"
	"go-product-service/config"
	"go-product-service/pkg/db"
)

func main() {
	db.Init(config.LoadConfig().DatabaseUrl)

	r := gin.Default()

	r.Run("localhost:8080")
}
