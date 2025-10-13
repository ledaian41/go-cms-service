package node_type_service

import (
	"slices"
	"time"

	"github.com/ledaian41/go-cms-service/pkg/helper/sql_helper"
	"github.com/ledaian41/go-cms-service/pkg/shared/dto"
	"github.com/ledaian41/go-cms-service/pkg/shared/utils"
	"github.com/ledaian41/go-cms-service/pkg/value_type"
	"gorm.io/gorm"
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
	db := s.db.Table(tid)
	db = db.Where("deleted_at IS NULL")

	var hasReference bool
	var joinSpec sql_helper.JoinSpec
	referenceView := option.GetReferenceViewKeys()
	if len(referenceView) > 0 {
		propertyTypes := s.FetchPropertyTypesByTid(tid)
		var referencePts []shared_dto.PropertyTypeDTO
		for _, pt := range propertyTypes {
			contain := slices.Contains(referenceView, pt.PID)
			reference := string(value_type.Reference)
			if contain && (pt.ValueType == reference || pt.ValueType == string(value_type.References)) {
				referencePts = append(referencePts, pt)
			}
		}
		hasReference = len(referencePts) > 0
		if hasReference {
			joinSpec = sql_helper.NewJoinSpec(tid, referencePts)
			query := sql_helper.QueryJoin(joinSpec)
			if len(query) > 0 {
				db.Select(sql_helper.BuildSelectFields(tid, joinSpec)).Joins(query)
			}
		}
	}

	searchQuery := option.GetSearchQuery()
	if len(searchQuery) > 0 {
		whereClause, values := sql_helper.BuildSearchConditions(searchQuery)
		if values != nil {
			db = db.Where(whereClause, values...)
		}
	}
	if len(option.SortBy) > 0 {
		db.Order(option.SortBy)
	}

	var total int64
	db.Count(&total)

	db.Limit(int(option.PageSize)).Offset(offset)
	if err := db.Find(&records).Error; err != nil {
		return nil, nil, err
	}

	pagination := &shared_dto.PaginationDTO{
		Page:     option.Page,
		PageSize: option.PageSize,
		Total:    total,
	}
	pagination.CalculateTotalPage()

	if hasReference {
		records = sql_helper.FormatJoinResponse(records, joinSpec)
	}
	return records, pagination, nil
}

func (s *NodeTypeService) FetchRecord(tid string, id string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := s.db.Table(tid).Where("id = ? AND deleted_at IS NULL", id).Find(&result).Error; err != nil {
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
	result := s.db.Table(tid).Where("id = ? AND deleted_at IS NULL", id).Updates(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (s *NodeTypeService) DeleteRecord(tid string, id string) error {
	// Soft delete: mark deleted_at and optionally deleted_by
	updates := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	// NOTE: if you have the actor/user id available at handler layer, pass it down and set deleted_by accordingly.
	// e.g., updates["deleted_by"] = actorID
	return s.db.Table(tid).Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error
}

func (s *NodeTypeService) RestoreRecord(tid string, id string) error {
	// Soft-restore: clear deleted_at and deleted_by only if the record is currently soft-deleted
	updates := map[string]interface{}{
		"deleted_at": nil,
	}
	result := s.db.Table(tid).Where("id = ? AND deleted_at IS NOT NULL", id).Updates(updates)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
