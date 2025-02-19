package helper_handler

import (
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
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *HelperHandler) FetchNodeType(c *gin.Context) {
	nodeTypes := h.nodeTypeService.FetchNodeTypes()
	c.JSON(http.StatusOK, gin.H{"items": nodeTypes})
}
