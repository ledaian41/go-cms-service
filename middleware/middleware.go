package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
	"github.com/ledaian41/go-cms-service/pkg/shared/interface"
	"net/http"
)

func CheckNodeTypeExist(nodeTypeService shared_interface.NodeTypeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		for i := range c.Params {
			if c.Params[i].Key == "typeId" {
				c.Params[i].Value = strcase.ToSnake(c.Params[i].Value)
				break
			}
		}

		typeId := c.Param("typeId")
		if !nodeTypeService.CheckNodeTypeExist(typeId) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("%s does not exist", typeId)})
			return
		}

		c.Set("typeId", typeId)
		c.Next()
	}
}
