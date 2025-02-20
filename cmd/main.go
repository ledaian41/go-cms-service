package main

import (
	"github.com/gin-gonic/gin"
	"go-product-service/config"
	"go-product-service/pkg/db"
	"go-product-service/pkg/helper/handler"
	nodeType_service "go-product-service/pkg/nodeType/service"
)

func main() {
	db := db.Init(config.LoadConfig().CachePath)
	nodeTypeService := nodeType_service.NewNodeTypeService(db)
	nodeTypeService.InitDatabase()

	r := gin.Default()

	helperHandler := helper_handler.NewHelperHandler(nodeTypeService)
	r.GET("helper/loadSchema", helperHandler.LoadSchema)
	r.GET("helper/nodeType", helperHandler.FetchNodeType)

	r.Run("localhost:8080")
}
