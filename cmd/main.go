package main

import (
	"github.com/gin-gonic/gin"
	"go-cms-service/config"
	"go-cms-service/middleware"
	"go-cms-service/pkg/db"
	"go-cms-service/pkg/helper/handler"
	"go-cms-service/pkg/nodeType/handler"
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

	nodeTypeHandler := nodeType_handler.NewNodeTypeHandler(nodeTypeService)
	r.GET("/:typeId", middleware.CheckNodeTypeExist(nodeTypeService), nodeTypeHandler.ListApi)
	r.GET("/:typeId/:id", middleware.CheckNodeTypeExist(nodeTypeService), nodeTypeHandler.ReadApi)
	r.POST("/:typeId", middleware.CheckNodeTypeExist(nodeTypeService), nodeTypeHandler.CreateApi)
	r.PATCH("/:typeId", middleware.CheckNodeTypeExist(nodeTypeService), nodeTypeHandler.UpdateApi)
	r.DELETE("/:typeId/:id", middleware.CheckNodeTypeExist(nodeTypeService), nodeTypeHandler.DeleteApi)

	r.Run("localhost:8080")
}
