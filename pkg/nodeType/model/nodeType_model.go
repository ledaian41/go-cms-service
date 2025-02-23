package nodeType_model

import (
	"go-cms-service/pkg/shared/dto"
	shared_utils "go-cms-service/pkg/shared/utils"
	"gorm.io/gorm"
)

type NodeType struct {
	gorm.Model
	ID            string          `gorm:"primaryKey;type:char(8)"`
	TID           string          `json:"tid" gorm:"uniqueIndex;column:tid"`
	PropertyTypes []*PropertyType `json:"propertyTypes" gorm:"foreignKey:NodeTypeRefer"`
}

func (n *NodeType) BeforeCreate(_ *gorm.DB) (err error) {
	n.ID = shared_utils.RandomID()
	return
}

type PropertyType struct {
	gorm.Model
	ID            string `gorm:"primaryKey;type:char(8)"`
	NodeTypeRefer string
	PID           string `json:"pid" gorm:"column:pid"`
	ValueType     string `json:"valueType"`
}

func (n *PropertyType) BeforeCreate(_ *gorm.DB) (err error) {
	n.ID = shared_utils.RandomID()
	return
}

func (n *NodeType) NodeTypeDTO() shared_dto.NodeTypeDTO {
	var propertyTypeDTOs []shared_dto.PropertyTypeDTO
	for _, prop := range n.PropertyTypes {
		propertyTypeDTOs = append(propertyTypeDTOs, shared_dto.PropertyTypeDTO{
			PID:       prop.PID,
			ValueType: prop.ValueType,
		})
	}
	return shared_dto.NodeTypeDTO{
		TID:           n.TID,
		PropertyTypes: propertyTypeDTOs,
	}
}
