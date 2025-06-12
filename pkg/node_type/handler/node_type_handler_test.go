package node_type_handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ledaian41/go-cms-service/pkg/shared/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockNodeTypeService struct {
	mock.Mock
}

func (m *MockNodeTypeService) FetchNodeTypes() *[]shared_dto.NodeTypeDTO {
	//TODO implement me
	panic("implement me")
}

func (m *MockNodeTypeService) FetchNodeType(tid string) shared_dto.NodeTypeDTO {
	panic("implement me")
}

func (m *MockNodeTypeService) LoadSchema(filePath string, ch chan<- string) {
	//TODO implement me
	panic("implement me")
}

func (m *MockNodeTypeService) DeleteNodeType(tid string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockNodeTypeService) CheckNodeTypeExist(tid string) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MockNodeTypeService) FetchRecords(tid string) ([]map[string]interface{}, error) {
	args := m.Called(tid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockNodeTypeService) FetchRecord(tid string, id string) (map[string]interface{}, error) {
	args := m.Called(tid, id)
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil
		}
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockNodeTypeService) CreateRecord(tid string, data map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(tid, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockNodeTypeService) UpdateRecord(tid string, id string, data map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(tid, id, data)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockNodeTypeService) DeleteRecord(tid string, id string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockNodeTypeService) PreprocessFile(nodeTypeDTO shared_dto.NodeTypeDTO, rawData map[string]interface{}) (map[string]interface{}, error) {
	panic("implement me")
}

func TestListApi_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNodeTypeService)
	mockData := []map[string]interface{}{
		{"id": 1, "name": "product A", "price": 200000},
		{"id": 2, "name": "product B", "price": 400000},
	}
	mockService.On("FetchRecords", "product").Return(&mockData, nil)

	handler := NewNodeTypeHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/product", nil)
	c.Params = gin.Params{
		{Key: "typeId", Value: "product"},
	}

	handler.ListApi(c)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedResponse := `[{"id":1,"name":"product A","price":200000},{"id":2,"name":"product B","price":400000}]`
	assert.JSONEq(t, expectedResponse, w.Body.String())

	mockService.AssertExpectations(t)
}

func TestListApi_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNodeTypeService)
	mockService.On("FetchRecords", "product").Return(nil, errors.New("record not found"))

	handler := NewNodeTypeHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/product", nil)
	c.Params = gin.Params{
		{Key: "typeId", Value: "product"},
	}

	handler.ListApi(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	assert.Contains(t, w.Body.String(), "record not found")

	mockService.AssertExpectations(t)
}

func TestReadApi_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNodeTypeService)
	mockData := map[string]interface{}{
		"id": 1, "name": "product A", "price": 200000,
	}
	mockService.On("FetchRecord", "product", "1").Return(&mockData, nil)

	handler := NewNodeTypeHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/product/1", nil)
	c.Params = gin.Params{
		{Key: "typeId", Value: "product"},
		{Key: "id", Value: "1"},
	}

	handler.ReadApi(c)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedResponse := `{"id":1,"name":"product A","price":200000}`
	assert.JSONEq(t, expectedResponse, w.Body.String())

	mockService.AssertExpectations(t)
}

func TestReadApi_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNodeTypeService)
	mockService.On("FetchRecord", "product", "1").Return(nil, errors.New("db error"))

	handler := NewNodeTypeHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/product/1", nil)
	c.Params = gin.Params{
		{Key: "typeId", Value: "product"},
		{Key: "id", Value: "1"},
	}

	handler.ReadApi(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "db error", w.Body.String())

	mockService.AssertExpectations(t)
}

func TestReadApi_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNodeTypeService)
	mockService.On("FetchRecord", "product", "1").Return(nil, nil)

	handler := NewNodeTypeHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/product/1", nil)
	c.Params = gin.Params{
		{Key: "typeId", Value: "product"},
		{Key: "id", Value: "1"},
	}

	handler.ReadApi(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, w.Body.String(), "not found")

	mockService.AssertExpectations(t)
}

func TestCreateApi_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNodeTypeService)

	requestBody := `{"name": "New Product", "price": 300000}`
	requestData := map[string]interface{}{"name": "New Product", "price": float64(300000)}
	createdData := map[string]interface{}{"id": 1, "name": "New Product", "price": 300000}

	mockService.On("CreateRecord", "product", requestData).Return(&createdData, nil)

	handler := NewNodeTypeHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "/product", strings.NewReader(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "typeId", Value: "product"},
	}

	handler.CreateApi(c)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedJSON := `{"id":1,"name":"New Product","price":300000}`
	assert.JSONEq(t, expectedJSON, w.Body.String())

	mockService.AssertExpectations(t)
}

func TestCreateApi_InvalidRequestData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNodeTypeService)
	handler := NewNodeTypeHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	invalidRequestBody := `{"name": "New Product", price: 300000}` // invalid body, missing "
	c.Request, _ = http.NewRequest(http.MethodPost, "/product", strings.NewReader(invalidRequestBody))
	c.Params = gin.Params{{Key: "typeId", Value: "product"}}

	handler.CreateApi(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid character")

	mockService.AssertExpectations(t)
}

func TestCreateApi_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNodeTypeService)
	requestBody := `{"name": "New Product", "price": 300000}`
	requestData := map[string]interface{}{"name": "New Product", "price": float64(300000)}
	mockService.On("CreateRecord", "product", requestData).Return(nil, errors.New("db error"))

	handler := NewNodeTypeHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/product", strings.NewReader(requestBody))
	c.Params = gin.Params{{Key: "typeId", Value: "product"}}

	handler.CreateApi(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "db error", w.Body.String())

	mockService.AssertExpectations(t)
}

func TestUpdateApi_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockNodeTypeService)
	requestBody := `{"name": "Updated Product", "price": 300000}`
	requestChangeData := map[string]interface{}{"name": "Updated Product", "price": float64(300000)}
	mockCurrentData := map[string]interface{}{"id": 1, "name": "New Product", "price": 200000}
	updatedData := map[string]interface{}{"id": 1, "name": "Updated Product", "price": 300000}
	mockService.On("FetchRecord", "product", "1").Return(&mockCurrentData, nil)
	mockService.On("UpdateRecord", "product", "1", requestChangeData).Return(&updatedData, nil)

	handler := NewNodeTypeHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPatch, "/product/1", strings.NewReader(requestBody))
	c.Params = gin.Params{{Key: "typeId", Value: "product"}, {Key: "id", Value: "1"}}

	handler.UpdateApi(c)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedJSON := `{"id":1,"name":"Updated Product","price":300000}`
	assert.JSONEq(t, expectedJSON, w.Body.String())
}
