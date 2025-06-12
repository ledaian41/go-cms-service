package helper_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ledaian41/go-cms-service/pkg/shared/interface"
	"log"
	"net/http"
)

type HelperHandler struct {
	nodeTypeService shared_interface.NodeTypeService
	helperService   shared_interface.HelperService
}

func NewHelperHandler(nodeTypeService shared_interface.NodeTypeService, helperService shared_interface.HelperService) *HelperHandler {
	return &HelperHandler{nodeTypeService: nodeTypeService, helperService: helperService}
}

func (h *HelperHandler) LoadSchema(c *gin.Context) {
	filePath := c.Query("filePath")
	messageCh := make(chan string)
	go h.nodeTypeService.LoadSchema(filePath, messageCh)
	// Server-Sent Events (SSE) - Config Header streaming message
	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Status(http.StatusOK)

	for tid := range messageCh {
		_, err := c.Writer.Write([]byte(fmt.Sprintf("🎉 Load nodeType: %s successfully!\n", tid)))
		if err != nil {
			log.Printf("❌ Error writing to response: %v", err)
			break
		}
		c.Writer.Flush()
	}
}

func (h *HelperHandler) LoadData(c *gin.Context) {
	filePath := c.Query("filePath")
	messageCh := make(chan string)
	go h.helperService.LoadJsonData(filePath, messageCh)

	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Status(http.StatusOK)

	for node := range messageCh {
		_, err := c.Writer.Write([]byte(fmt.Sprintf("🎉 Load %s successfully!\n", node)))
		if err != nil {
			log.Printf("❌ Error writing to response: %v", err)
			break
		}
		c.Writer.Flush()
	}
}

func (h *HelperHandler) FetchNodeType(c *gin.Context) {
	nodeTypes := h.nodeTypeService.FetchNodeTypes()
	c.JSON(http.StatusOK, gin.H{"items": nodeTypes})
}

func (h *HelperHandler) DeleteNodeType(c *gin.Context) {
	tid := c.Query("typeId")
	success, err := h.nodeTypeService.DeleteNodeType(tid)
	if !success {
		c.String(http.StatusOK, fmt.Sprintf("🎉 Delete nodeType: %s - %v", tid, err))
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("🎉 Delete nodeType: %s successfully!", tid))
}
