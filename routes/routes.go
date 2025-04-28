package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go-cms-service/config"
	_ "go-cms-service/docs"
	"go-cms-service/middleware"
	"go-cms-service/pkg/file/handler"
	"go-cms-service/pkg/file/service"
	"go-cms-service/pkg/helper/handler"
	"go-cms-service/pkg/helper/service"
	"go-cms-service/pkg/nodeType/handler"
	"go-cms-service/pkg/nodeType/service"
	"gorm.io/gorm"
	"net/http"
)

func InitRoutes(db *gorm.DB, redis *config.RedisClient) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fileService := file_service.NewFileService()
	nodeTypeService := nodeType_service.NewNodeTypeService(db, fileService)
	nodeTypeService.InitDatabase()

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

	fileHandler := handler.NewFileHandler(fileService)
	r.GET("/file/*path", fileHandler.ReadFile)

	return r
}
