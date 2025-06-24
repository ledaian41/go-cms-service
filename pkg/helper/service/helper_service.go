package helper_service

import (
	"encoding/json"
	"fmt"
	"github.com/ledaian41/go-cms-service/pkg/helper/sql_helper"
	"github.com/ledaian41/go-cms-service/pkg/shared/utils"
	"gorm.io/gorm"
	"io/ioutil"
	"sync"
)

type HelperService struct {
	db               *gorm.DB
	tableColumnCache sync.Map
}

func NewHelperService(db *gorm.DB) *HelperService {
	return &HelperService{db: db, tableColumnCache: sync.Map{}}
}

func (s *HelperService) LoadJsonData(path string, ch chan<- string) {
	if shared_utils.IsJsonPath(path) {
		content := loadJsonFile(path)
		s.loadJsonToDB(content, ch)
		s.tableColumnCache.Clear()
	}

}

func (s *HelperService) loadJsonToDB(content []map[string]interface{}, ch chan<- string) {
	for _, item := range content {
		if typeId, ok := item["type_id"].(string); ok {
			if !s.db.Migrator().HasTable(typeId) {
				continue
			}

			id, ok := item["id"].(string)
			if ok {
				if s.checkRecordExist(typeId, id) {
					recordId := s.updateRecord(typeId, item)
					ch <- fmt.Sprintf("%s::%s", typeId, recordId)
					continue
				}
			}

			recordId := s.createNewRecord(typeId, item)
			if recordId != "" {
				ch <- fmt.Sprintf("%s::%s", typeId, recordId)
			}
		}
	}
	close(ch)
}

func (s *HelperService) checkRecordExist(tid string, id string) bool {
	var count int64
	if err := s.db.Table(tid).Where("id = ?", id).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (s *HelperService) createNewRecord(tid string, record map[string]interface{}) string {
	columns := s.getTableColumns(tid)
	if columns == nil {
		return ""
	}

	if id, ok := record["id"].(string); !ok || len(id) == 0 {
		record["id"] = sql_helper.GenerateID()
	}

	validRecord := make(map[string]interface{})
	for col, val := range record {
		if _, exists := columns[col]; exists {
			validRecord[col] = val
		}
	}

	if len(validRecord) == 0 {
		return ""
	}

	err := s.db.Table(tid).Create(&validRecord).Error
	if err != nil {
		return ""
	}
	return record["id"].(string)
}

func (s *HelperService) updateRecord(tid string, record map[string]interface{}) string {
	columns := s.getTableColumns(tid)
	if columns == nil {
		return ""
	}

	validRecord := make(map[string]interface{})
	for col, val := range record {
		if _, exists := columns[col]; exists {
			validRecord[col] = val
		}
	}

	if len(validRecord) == 0 {
		return ""
	}

	recordId := record["id"].(string)
	s.db.Table(tid).Where("id = ?", recordId).Updates(&validRecord)
	return recordId
}

func (s *HelperService) getTableColumns(tid string) map[string]bool {
	if cols, ok := s.tableColumnCache.Load(tid); ok {
		return cols.(map[string]bool)
	}

	rows, err := s.db.Raw(sql_helper.QueryTableColumns(tid)).Rows()
	if err != nil {
		return nil
	}

	columns := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		columns[name] = true
	}

	s.tableColumnCache.Store(tid, columns)
	return columns
}

func loadJsonFile(path string) []map[string]interface{} {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	var raw json.RawMessage
	if err := json.Unmarshal(file, &raw); err != nil {
		return nil
	}

	if len(raw) > 0 {
		if raw[0] == '[' {
			var data []map[string]interface{}
			json.Unmarshal(raw, &data)
			return data
		}

		if raw[0] == '{' {
			var data map[string]interface{}
			json.Unmarshal(raw, &data)
			result := make([]map[string]interface{}, 1)
			result = append(result, data)
			return result
		}
	}
	return nil
}
