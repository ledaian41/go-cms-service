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

// ListApi godoc
// @Summary List all nodes by type
// @Description Get all nodes of a specific type
// @Tags NodeType
// @Accept json
// @Produce json
// @Param typeId path string true "Type ID"
// @Success 200
// @Failure 400
// @Router /{typeId} [get]
func (n *NodeType) ListApi(c *gin.Context) {
	result, err := n.nodeTypeService.FetchRecords(c.Param("typeId"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// ReadApi godoc
// @Summary	Get node details
// @Description Get detailed information of a specific node
// @Tags	NodeType
// @Accept	json
// @Produce json
// @Param typeId path string true "Type ID"
// @Param id path string true "Node ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Router /{typeId}/{id} [get]
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

// CreateApi godoc
// @Summary Create new node
// @Description Create a new node with form data
// @Tags NodeType
// @Accept multipart/form-data
// @Produce json
// @Param typeId path string true "Type ID"
// @Param title formData string true "Node title"
// @Param content formData string false "Node content"
// @Param image formData file false "Image file"
// @Success 200
// @Failure 400
// @Router /{typeId} [post]
func (n *NodeType) CreateApi(c *gin.Context) {
	typeId := c.Param("typeId")

	rawData := make(map[string]interface{})
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for key, values := range form.Value {
		if len(values) > 0 {
			rawData[key] = values[0]
		}
	}
	for key, files := range form.File {
		if len(files) > 0 {
			rawData[key] = files[0]
		}
	}

	parsedData, err := n.nodeTypeService.PreprocessFile(n.nodeTypeService.FetchNodeType(typeId), rawData)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	newNode, err := n.nodeTypeService.CreateRecord(typeId, parsedData)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, newNode)
}

// UpdateApi godoc
// @Summary Update existing node
// @Description Update node information
// @Tags NodeType
// @Accept multipart/form-data
// @Produce json
// @Param typeId path string true "Type ID"
// @Param id path string true "Node ID"
// @Param title formData string false "Node title"
// @Param content formData string false "Node content"
// @Param image formData file false "Image file"
// @Success 200
// @Failure 400
// @Failure 404
// @Router /{typeId}/{id} [put]
func (n *NodeType) UpdateApi(c *gin.Context) {
	typeId := c.Param("typeId")
	id := c.Param("id")

	rawData := make(map[string]interface{})
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for key, values := range form.Value {
		if len(values) > 0 {
			rawData[key] = values[0]
		}
	}
	for key, files := range form.File {
		if len(files) > 0 {
			rawData[key] = files[0]
		}
	}

	record, err := n.nodeTypeService.FetchRecord(typeId, id)
	if err != nil || record == nil {
		c.String(http.StatusNotFound, "not found")
		return
	}

	parsedData, err := n.nodeTypeService.PreprocessFile(n.nodeTypeService.FetchNodeType(typeId), rawData)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	updateNode, err := n.nodeTypeService.UpdateRecord(typeId, id, parsedData)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, updateNode)
}

// DeleteApi godoc
// @Summary Delete node
// @Description Permanently delete a specific node
// @Tags NodeType
// @Accept json
// @Produce json
// @Param typeId path string true "Type ID"
// @Param id path string true "Node ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Router /{typeId}/{id} [delete]
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
