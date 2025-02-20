package helper_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-product-service/pkg/shared/interface"
	"net/http"
)

type HelperHandler struct {
	nodeTypeService shared_interface.NodeTypeServiceInterface
}

func NewHelperHandler(nodeTypeService shared_interface.NodeTypeServiceInterface) *HelperHandler {
	return &HelperHandler{nodeTypeService: nodeTypeService}
}

func (h *HelperHandler) LoadSchema(c *gin.Context) {
	filePath := c.Query("filePath")
	tid, err := h.nodeTypeService.LoadSchema(filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	c.String(http.StatusOK, fmt.Sprintf("Load %s schema successfully!", tid))
}

func (h *HelperHandler) FetchNodeType(c *gin.Context) {
	nodeTypes := h.nodeTypeService.FetchNodeTypes()
	c.JSON(http.StatusOK, gin.H{"items": nodeTypes})
}
