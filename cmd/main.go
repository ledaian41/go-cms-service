package main

import (
	"github.com/gin-gonic/gin"
	"go-cms-service/config"
	"go-cms-service/pkg/db"
	"go-cms-service/pkg/helper/handler"
	"go-cms-service/pkg/nodeType/service"
)

func main() {
	db := db.Init(config.LoadConfig().CachePath)
	nodeTypeService := nodeType_service.NewNodeTypeService(db)
	nodeTypeService.InitDatabase()

	r := gin.Default()

	helperHandler := helper_handler.NewHelperHandler(nodeTypeService)
	r.GET("helper/loadSchema", helperHandler.LoadSchema)
	r.GET("helper/nodeType", helperHandler.FetchNodeType)
	r.GET("helper/nodeType/delete", helperHandler.DeleteNodeType)

	r.Run("localhost:8080")
}
