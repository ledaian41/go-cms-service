package helper_service

import (
	"errors"
	"fmt"
	"github.com/ledaian41/go-cms-service/pkg/helper/sql_helper"
	"github.com/ledaian41/go-cms-service/pkg/node_type/model"
	"github.com/ledaian41/go-cms-service/pkg/node_type/utils"
	"github.com/ledaian41/go-cms-service/pkg/shared/utils"
	"github.com/ledaian41/go-cms-service/pkg/value_type"
	"gorm.io/gorm"
	"log"
	"strings"
	"sync"
)

func (s *HelperService) LoadSchema(path string, ch chan<- string) {
	if shared_utils.IsJsonPath(path) {
		s.loadSchemaFile(path, ch)
		return
	}
	if shared_utils.IsDirectory(path) {
		s.loadSchemaDirectory(path, ch)
		return
	}
}

func (s *HelperService) loadSchemaFile(path string, ch chan<- string) {
	nodeType, err := nodeType_utils.ReadSchemaJson(path)
	if err != nil {
		log.Printf("❌ Failed at LoadSchema: %v", err)
		close(ch)
		return
	}
	s.loadNodeTypeToDB(nodeType, ch)
	close(ch)
}

func (s *HelperService) loadSchemaDirectory(path string, ch chan<- string) {
	nodeTypes, err := nodeType_utils.ReadSchemasFromDir(path)
	if err != nil {
		log.Printf("❌ Failed at LoadSchema: %v", err)
		close(ch)
		return
	}

	var wg sync.WaitGroup
	for _, nodeType := range nodeTypes {
		wg.Add(1)
		go func(nodeType *node_type_model.NodeType) {
			defer wg.Done()
			s.loadNodeTypeToDB(nodeType, ch)
		}(nodeType)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
}

func (s *HelperService) loadNodeTypeToDB(nodeType *node_type_model.NodeType, ch chan<- string) {
	var existing node_type_model.NodeType
	if err := s.db.Preload("PropertyTypes").Where("tid = ?", nodeType.TID).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tid, _ := s.createNewNodeType(nodeType)
			ch <- tid
		} else {
			log.Printf("❌ Error loading record: %v", err)
		}
		return
	}
	tid, _ := s.updateNodeType(&existing, nodeType)
	ch <- tid
}

func (s *HelperService) createNewNodeType(nodeType *node_type_model.NodeType) (string, error) {
	if err := s.db.Exec(sql_helper.QueryCreateNewTable(nodeType)).Error; err != nil {
		log.Printf("❌ Failed at create Table: %v", err)
		return nodeType.TID, err
	}
	if err := s.db.Create(&nodeType).Error; err != nil {
		log.Printf("❌ Failed at save NodeType: %v", err)
		return nodeType.TID, err
	}
	log.Printf("🎉 Helper - Load new %s schema successfully!", nodeType.TID)
	return nodeType.TID, nil
}

func (s *HelperService) deleteColumn(tid, pid string) error {
	log.Printf("Delete column: %s in table %s", pid, tid)
	return s.db.Exec(sql_helper.QueryDeleteColumnFromTable(tid, pid)).Error
}

func (s *HelperService) updateNodeType(existing *node_type_model.NodeType, newNodeType *node_type_model.NodeType) (string, error) {
	var currentPTs []*node_type_model.PropertyType
	if err := s.db.Model(&existing).Association("PropertyTypes").Find(&currentPTs); err != nil {
		log.Printf("❌ Failed at query PropertyTypes: %v", err)
		return newNodeType.TID, nil
	}

	currentMap := make(map[string]*node_type_model.PropertyType)
	for _, pt := range currentPTs {
		currentMap[pt.PID] = pt
	}

	newMap := make(map[string]*node_type_model.PropertyType)
	var toCreate []*node_type_model.PropertyType
	for _, pt := range newNodeType.PropertyTypes {
		pt.NodeTypeRefer = existing.ID
		if currentMap[pt.PID] != nil {
			newMap[pt.PID] = pt
		} else {
			toCreate = append(toCreate, pt)
		}
	}

	for pid, pt := range currentMap {
		if newPT, ok := newMap[pid]; ok {
			if value_type.MapValueTypeToSQL(pt) != value_type.MapValueTypeToSQL(newPT) {
				if err := s.deleteColumn(newNodeType.TID, pid); err != nil {
					if !strings.Contains(err.Error(), "no such column") {
						log.Printf("❌ Error delete column %s: %v\n", pt.PID, err)
						return newNodeType.TID, nil
					}
				}
				if err := s.db.Unscoped().Delete(pt).Error; err != nil {
					log.Printf("❌ Failed to delete PropertyType (pid=%s): %v", pid, err)
					return newNodeType.TID, nil
				}
				toCreate = append(toCreate, newPT)
			} else {
				pt.ValueType = newPT.ValueType
				if err := s.db.Save(pt).Error; err != nil {
					log.Printf("❌ Failed to update PropertyType (pid=%s): %v", pid, err)
					return newNodeType.TID, nil
				}
			}
		} else {
			if err := s.deleteColumn(newNodeType.TID, pid); err != nil {
				if !strings.Contains(err.Error(), "no such column") {
					log.Printf("❌ Error delete column %s: %v\n", pt.PID, err)
					return newNodeType.TID, nil
				}
			}
			if err := s.db.Unscoped().Delete(pt).Error; err != nil {
				log.Printf("❌ Failed to delete PropertyType (pid=%s): %v", pid, err)
				return newNodeType.TID, nil
			}
		}
	}

	for _, pt := range toCreate {
		sql := sql_helper.QueryAddColumnToTable(newNodeType.TID, pt)
		if len(sql) == 0 {
			continue
		}
		if err := s.db.Exec(sql).Error; err != nil {
			log.Printf("❌ Failed at AutoMigrate: %v", err)
			return newNodeType.TID, nil
		}
		fmt.Println("create new PropertyType", pt.PID)
		if err := s.db.Create(pt).Error; err != nil {
			log.Printf("❌ Failed to create new PropertyType: %v", err)
			return newNodeType.TID, nil
		}
	}

	newNodeType.ID = existing.ID
	if err := s.db.Omit("PropertyTypes").Save(&newNodeType).Error; err != nil {
		log.Printf("❌ Failed at update NodeType(%s): %v", newNodeType.TID, err)
		return newNodeType.TID, err
	}
	log.Printf("🎉 Helper - Load %s schema successfully!", newNodeType.TID)
	return newNodeType.TID, nil
}
