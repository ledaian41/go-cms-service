package utils

import (
	"encoding/json"
	"go-product-service/pkg/nodeType/model"
	"io/ioutil"
)

func loadSchema(path string) (*model.NodeType, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var schema model.NodeType
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, err
	}
	return &schema, nil
}
