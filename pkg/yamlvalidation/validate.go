package yamlvalidation

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"
)

var yamlToJson = yaml.YAMLToJSON
var validateJson = gojsonschema.Validate

func Validate(values map[string]interface{}, schemaJSON []byte) (reterr error) {
	valuesData, err := yaml.Marshal(values)
	if err != nil {
		return err
	}

	valuesJSON, err := yamlToJson(valuesData)
	if err != nil {
		return err
	}

	if bytes.Equal(valuesJSON, []byte("null")) {
		valuesJSON = []byte("{}")
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaJSON)
	valuesLoader := gojsonschema.NewBytesLoader(valuesJSON)

	result, err := validateJson(schemaLoader, valuesLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		var sb strings.Builder
		for _, desc := range result.Errors() {
			sb.WriteString(fmt.Sprintf("- %s\n", desc))
		}
		return errors.New(sb.String())
	}

	return nil
}
