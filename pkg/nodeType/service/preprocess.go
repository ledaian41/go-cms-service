package nodeType_service

import (
	"go-cms-service/pkg/shared/dto"
	"mime/multipart"
	"os"
)

func (s *NodeTypeService) PreprocessData(nodeTypeDTO shared_dto.NodeTypeDTO, rawData map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range rawData {
		result[k] = v
	}

	propertyTypes := nodeTypeDTO.PropertyTypes
	for _, propertyType := range propertyTypes {
		if propertyType.ValueType != "FILE" {
			continue
		}

		fileHeader, ok := rawData[propertyType.PID].(*multipart.FileHeader)
		if !ok {
			continue
		}

		fileInfo, err := s.fileService.SaveFile(fileHeader, "./cache/files/")
		if err != nil {
			if (fileInfo != nil) && (fileInfo.SavedPath != "") {
				os.Remove(fileInfo.SavedPath)
			}
			continue
		}

		result[propertyType.PID] = fileInfo.SavedPath
	}
	return result
}
