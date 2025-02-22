package nodeType_utils

import (
	"encoding/json"
	"go-product-service/pkg/nodeType/model"
	"io/ioutil"
)

func ReadSchemaJson(path string) (*nodeType_model.NodeType, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var schema nodeType_model.NodeType
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, err
	}
	return &schema, nil
}
