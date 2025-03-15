package nodeType_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-cms-service/pkg/shared/interface"
	"net/http"
)

type NodeType struct {
	nodeTypeService shared_interface.NodeTypeServiceInterface
}

func NewNodeTypeHandler(nodeTypeService shared_interface.NodeTypeServiceInterface) *NodeType {
	return &NodeType{nodeTypeService: nodeTypeService}
}

func (n *NodeType) checkTypeId(typeId string) error {
	if n.nodeTypeService.CheckNodeTypeExist(typeId) {
		return nil
	}
	return fmt.Errorf("%s does not exist", typeId)
}

func (n *NodeType) ListApi(c *gin.Context) {
	result, err := n.nodeTypeService.FetchRecords(c.Param("typeId"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (n *NodeType) ReadApi(c *gin.Context) {
	result, err := n.nodeTypeService.FetchRecord(c.Param("typeId"), c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if *result == nil {
		c.String(http.StatusNotFound, "not found")
		return
	}
	c.JSON(http.StatusOK, result)
}

func (n *NodeType) CreateApi(c *gin.Context) {
	typeId := c.Param("typeId")

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newNode, err := n.nodeTypeService.CreateRecord(typeId, data)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, newNode)
}

func (n *NodeType) UpdateApi(c *gin.Context) {
	typeId := c.Param("typeId")
	c.JSON(http.StatusOK, gin.H{"items": "Update Api", "typeId": typeId})
}

func (n *NodeType) DeleteApi(c *gin.Context) {
	typeId := c.Param("typeId")
	c.JSON(http.StatusOK, gin.H{"items": "Delete Api", "typeId": typeId})
}
