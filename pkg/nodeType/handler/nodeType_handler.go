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
	if result == nil {
		c.String(http.StatusNotFound, "not found")
		return
	}
	c.JSON(http.StatusOK, result)
}

func (n *NodeType) CreateApi(c *gin.Context) {
	typeId := c.Param("typeId")

	data := make(map[string]interface{})
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for key, values := range form.Value {
		if len(values) > 0 {
			data[key] = values[0]
		}
	}
	for key, files := range form.File {
		if len(files) > 0 {
			data[key] = files[0]
		}
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
	id := c.Param("id")

	data := make(map[string]interface{})
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for key, values := range form.Value {
		if len(values) > 0 {
			data[key] = values[0]
		}
	}
	for key, files := range form.File {
		if len(files) > 0 {
			data[key] = files[0]
		}
	}

	record, err := n.nodeTypeService.FetchRecord(typeId, id)
	if err != nil || record == nil {
		c.String(http.StatusNotFound, "not found")
		return
	}

	updateNode, err := n.nodeTypeService.UpdateRecord(typeId, id, data)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, updateNode)
}

func (n *NodeType) DeleteApi(c *gin.Context) {
	typeId := c.Param("typeId")
	id := c.Param("id")

	record, err := n.nodeTypeService.FetchRecord(typeId, id)
	if err != nil || record == nil {
		c.String(http.StatusNotFound, "not found")
		return
	}

	if err = n.nodeTypeService.DeleteRecord(typeId, id); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
