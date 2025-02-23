package nodeType_model

import (
	"go-product-service/pkg/shared/dto"
	"gorm.io/gorm"
)

type NodeType struct {
	gorm.Model
	TID           string          `json:"tid" gorm:"uniqueIndex;column:tid"`
	PropertyTypes []*PropertyType `json:"propertyTypes" gorm:"foreignKey:NodeTypeRefer"`
}

type PropertyType struct {
	gorm.Model
	NodeTypeRefer uint
	PID           string `json:"pid" gorm:"column:pid"`
	ValueType     string `json:"valueType"`
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
