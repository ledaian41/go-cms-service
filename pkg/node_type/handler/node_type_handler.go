package node_type_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
	"github.com/ledaian41/go-cms-service/pkg/shared/interface"
	"github.com/ledaian41/go-cms-service/pkg/shared/utils"
)

type NodeType struct {
	nodeTypeService shared_interface.NodeTypeService
}

func NewNodeTypeHandler(nodeTypeService shared_interface.NodeTypeService) *NodeType {
	return &NodeType{nodeTypeService: nodeTypeService}
}

// ListApi godoc
// @Summary List nodes by type
// @Description Get nodes of a specific type with pagination, sorting, and flexible filter syntax.
// @Description \n
// @Description **Filtering syntax** (all remaining URL query params are interpreted as filters):
// @Description - Pattern: `{field}_{operator}={value}`
// @Description - Supported operators: `equal`, `include`, `in`, `from`, `to`, `fromto`
// @Description - Semantics:
// @Description   * `equal`: exact match (e.g. `status_equal=published`)
// @Description   * `include`: substring/contains (e.g. `title_include=hello`)
// @Description   * `in`: membership list, comma-separated (e.g. `type_in=article,page`)
// @Description   * `from`: lower bound (>=), typically for dates/numbers (e.g. `createdAt_from=2025-01-01T00:00:00Z`)
// @Description   * `to`: upper bound (<=) (e.g. `createdAt_to=2025-12-31T23:59:59Z`)
// @Description   * `fromto`: range (e.g. `price_fromto=10,100`)
// @Description - Examples: `GET /{typeId}?title_include=guide&status_in=draft,published&createdAt_from=2025-01-01T00:00:00Z`
// @Description \n
// @Description **Sorting syntax**
// @Description - Pattern: `<field> <asc|desc>`; default direction is `asc` if omitted (e.g., `createdAt` == `createdAt asc`)
// @Description - Multiple fields: separate by comma, evaluated left-to-right (e.g., `name desc,age asc`)
// @Description - URL encoding: encode spaces as `%20` or `+` (e.g., `name%20desc,age%20asc`)
// @Description - Examples: `GET /{typeId}?sort=createdAt%20desc,id`, `GET /{typeId}?sort=name%20desc,age%20asc`
// @Tags NodeType
// @Accept json
// @Produce json
// @Param typeId path string true "Type ID"
// @Param page query int false "Page number (1-based)" default(1) minimum(1)
// @Param pageSize query int false "Items per page (1-1000)" default(10) minimum(1) maximum(1000)
// @Param sort query string false "Sort expression: `<field> <asc|desc>`, multiple fields separated by comma. Example: `name desc,age asc`. Default direction is `asc` if omitted. Use `%20` (or `+`) to encode spaces in URLs: `name%20desc,age%20asc`"
// @Param filter query string false "Dynamic filters: `{field}_{operator}={value}`. Operators: `equal|include|in|from|to|fromto`. Example: `name_equal=ABC&age_from=20`"
// @Param referenceView query string false "true or field name to fetch related records"
// @Success 200 {object} map[string]interface{} "{ items: [...], pagination: { page, pageSize, total, hasNext, nextCursor? } }"
// @Failure 400 {string} string "bad request"
// @Router /{typeId} [get]
func (n *NodeType) ListApi(c *gin.Context) {
	typeId := strcase.ToSnake(c.Param("typeId"))
	records, pagination, err := n.nodeTypeService.FetchRecords(typeId, shared_utils.QueryOption{
		TypeId:        typeId,
		ReferenceView: c.Query("referenceView"),
		PageSize:      int8(shared_utils.ParseInt(c.Query("pageSize"))),
		Page:          int32(shared_utils.ParseInt(c.Query("page"))),
		SortBy:        c.Query("sort"),
		Query:         c.Request.URL.Query(),
	})
	for _, record := range records {
		n.nodeTypeService.ProcessFilePath(record)
	}
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      records,
		"pagination": pagination,
	})
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
	n.nodeTypeService.ProcessFilePath(result)
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
