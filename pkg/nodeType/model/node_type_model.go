package nodeType_model

import (
	"go-cms-service/pkg/shared/dto"
	"go-cms-service/pkg/shared/utils"
	"gorm.io/gorm"
)

type NodeType struct {
	gorm.Model
	ID            string          `gorm:"primaryKey;type:char(8)"`
	TID           string          `json:"tid" gorm:"uniqueIndex;column:tid"`
	PropertyTypes []*PropertyType `json:"propertyTypes" gorm:"foreignKey:NodeTypeRefer"`
}

func (n *NodeType) BeforeCreate(_ *gorm.DB) (err error) {
	n.ID = shared_utils.RandomID(4)
	return
}

type PropertyType struct {
	gorm.Model
	ID             string `gorm:"primaryKey;type:char(8)"`
	NodeTypeRefer  string
	PID            string `json:"pid" gorm:"column:pid"`
	ValueType      string `json:"valueType"`
	ReferenceType  string `json:"referenceType"`
	ReferenceValue string `json:"referenceValue"`
}

func (pt *PropertyType) BeforeCreate(_ *gorm.DB) (err error) {
	pt.ID = shared_utils.RandomID(4)
	return
}

func (pt *PropertyType) PropertyTypeDTO() shared_dto.PropertyTypeDTO {
	return shared_dto.PropertyTypeDTO{
		PID:            pt.PID,
		ValueType:      pt.ValueType,
		ReferenceType:  pt.ReferenceType,
		ReferenceValue: pt.ReferenceValue,
	}
}

func (n *NodeType) NodeTypeDTO() shared_dto.NodeTypeDTO {
	var propertyTypeDTOs []shared_dto.PropertyTypeDTO
	for _, pt := range n.PropertyTypes {
		propertyTypeDTOs = append(propertyTypeDTOs, pt.PropertyTypeDTO())
	}
	return shared_dto.NodeTypeDTO{
		TID:           n.TID,
		PropertyTypes: propertyTypeDTOs,
	}
}
