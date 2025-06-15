package node_type_service

import (
	"github.com/ledaian41/go-cms-service/pkg/node_type/sql_helper"
	"github.com/ledaian41/go-cms-service/pkg/shared/dto"
	"github.com/ledaian41/go-cms-service/pkg/shared/utils"
)

func (s *NodeTypeService) FetchRecords(tid string, option shared_utils.QueryOption) ([]map[string]interface{}, *shared_dto.PaginationDTO, error) {
	var records []map[string]interface{}
	if option.Page < 1 {
		option.Page = 1 // default page number
	}

	if option.PageSize < 1 || option.PageSize > 100 {
		option.PageSize = 10 // default page size
	}
	offset := (int(option.Page) - 1) * int(option.PageSize)
	if err := s.db.Table(tid).Limit(int(option.PageSize)).Offset(offset).Find(&records).Error; err != nil {
		return nil, nil, err
	}

	var total int64
	if err := s.db.Table(tid).Count(&total).Error; err != nil {
		return nil, nil, err
	}
	pagination := &shared_dto.PaginationDTO{
		Page:     option.Page,
		PageSize: option.PageSize,
		Total:    total,
	}
	pagination.CalculateTotalPage()

	return records, pagination, nil
}

func (s *NodeTypeService) FetchRecord(tid string, id string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := s.db.Table(tid).Find(&result, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result, nil
}

func (s *NodeTypeService) CreateRecord(tid string, data map[string]interface{}) (map[string]interface{}, error) {
	data["id"] = sql_helper.GenerateID()
	if result := s.db.Table(tid).Create(&data); result.Error != nil {
		return data, result.Error
	}
	delete(data, "@id")
	return data, nil
}

func (s *NodeTypeService) UpdateRecord(tid string, id string, data map[string]interface{}) (map[string]interface{}, error) {
	delete(data, "id")
	result := s.db.Table(tid).Where("id = ?", id).Updates(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (s *NodeTypeService) DeleteRecord(tid string, id string) error {
	return s.db.Table(tid).Where("id = ?", id).Delete(nil).Error
}
