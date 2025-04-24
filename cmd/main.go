package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go-cms-service/config"
	_ "go-cms-service/docs"
	"go-cms-service/middleware"
	"go-cms-service/pkg/db"
	"go-cms-service/pkg/file/service"
	"go-cms-service/pkg/helper/handler"
	"go-cms-service/pkg/helper/service"
	"go-cms-service/pkg/nodeType/handler"
	"go-cms-service/pkg/nodeType/service"
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
	db := db.Init(fmt.Sprintf("%s/cache.sqlite", config.LoadConfig().CachePath))
	fileService := file_service.NewFileService()
	nodeTypeService := nodeType_service.NewNodeTypeService(db, fileService)
	nodeTypeService.InitDatabase()

	r := gin.Default()

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	helperService := helper_service.NewHelperService(db)
	helperHandler := helper_handler.NewHelperHandler(nodeTypeService, helperService)
	r.GET("helper/loadSchema", helperHandler.LoadSchema)
	r.GET("helper/loadData", helperHandler.LoadData)
	r.GET("helper/nodeType", helperHandler.FetchNodeType)
	r.GET("helper/nodeType/delete", helperHandler.DeleteNodeType)

	nodeTypeHandler := nodeType_handler.NewNodeTypeHandler(nodeTypeService)
	r.GET("/:typeId", middleware.CheckNodeTypeExist(nodeTypeService), nodeTypeHandler.ListApi)
	r.GET("/:typeId/:id", middleware.CheckNodeTypeExist(nodeTypeService), nodeTypeHandler.ReadApi)
	r.POST("/:typeId", middleware.CheckNodeTypeExist(nodeTypeService), nodeTypeHandler.CreateApi)
	r.PATCH("/:typeId/:id", middleware.CheckNodeTypeExist(nodeTypeService), nodeTypeHandler.UpdateApi)
	r.DELETE("/:typeId/:id", middleware.CheckNodeTypeExist(nodeTypeService), nodeTypeHandler.DeleteApi)

	r.Run("localhost:8080")
}
