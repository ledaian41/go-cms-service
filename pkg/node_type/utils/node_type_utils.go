package nodeType_utils

import (
	"encoding/json"
	"fmt"
	"github.com/ledaian41/go-cms-service/pkg/node_type/model"
	"io/ioutil"
	"log"
	"path/filepath"
)

func ReadSchemaJson(path string) (*node_type_model.NodeType, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var schema node_type_model.NodeType
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, err
	}
	return &schema, nil
}

func ReadSchemasFromDir(path string) ([]*node_type_model.NodeType, error) {
	var schemas []*node_type_model.NodeType
	pattern := filepath.Join(path, "*.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to glob files in directory: %s", path, err)
	}

	if len(files) == 0 {
		log.Printf("No files found in directory: %s", path)
		return nil, nil
	}

	for _, file := range files {
		schema, err := ReadSchemaJson(file)
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, schema)
	}
	return schemas, nil
}
