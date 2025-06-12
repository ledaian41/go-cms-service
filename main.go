package main

import (
	"go-cms-service/config"
	"go-cms-service/pkg/db"
	"go-cms-service/routes"
	"log"
	"os"
)

// @title Go CMS API
// @version 1.0
// @description A CMS API service.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.LoadConfig()
	redisClient := config.InitRedisClient()
	db := db.Init(config.Env.DbHost, config.Env.DbUser, config.Env.DbPwd)
	r := routes.InitRoutes(db, redisClient)
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	r.Run(":" + port)
}
