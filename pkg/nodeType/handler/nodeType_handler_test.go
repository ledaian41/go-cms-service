package nodeType_handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-cms-service/pkg/shared/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockNodeTypeService struct {
	mock.Mock
}

func (m *MockNodeTypeService) FetchNodeTypes() *[]shared_dto.NodeTypeDTO {
	//TODO implement me
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

func (m *MockNodeTypeService) FetchRecords(tid string) (*[]map[string]interface{}, error) {
	args := m.Called(tid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*[]map[string]interface{}), args.Error(1)
}

func (m *MockNodeTypeService) FetchRecord(tid string, id string) (*map[string]interface{}, error) {
	args := m.Called(tid, id)
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil
		}
		return nil, args.Error(1)
	}
	return args.Get(0).(*map[string]interface{}), args.Error(1)
}

func (m *MockNodeTypeService) CreateRecord(tid string, data map[string]interface{}) (*map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockNodeTypeService) UpdateRecord(tid string, id string, data map[string]interface{}) (*map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockNodeTypeService) DeleteRecord(tid string, id string) error {
	//TODO implement me
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

	c.Params = append(c.Params, gin.Param{Key: "typeId", Value: "product"})

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
	c.Params = gin.Params{{Key: "typeId", Value: "product"}}

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
	c.Params = append(c.Params, gin.Param{Key: "typeId", Value: "product"})
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

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
	c.Params = append(c.Params, gin.Param{Key: "typeId", Value: "product"})
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.ReadApi(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), "db error")

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
	c.Params = append(c.Params, gin.Param{Key: "typeId", Value: "product"})
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.ReadApi(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	//assert.Equal(t, w.Body.String(), "not found")

	mockService.AssertExpectations(t)
}
