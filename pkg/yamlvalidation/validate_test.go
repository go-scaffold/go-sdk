package yamlvalidation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
)

func TestValidateAgainstSchema(t *testing.T) {
	type mocks struct {
		yamlToJsonErr           error
		jsonSchemaValidationErr error
	}
	type args struct {
		values     map[string]interface{}
		schemaJSON []byte
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		wantErr string
	}{
		{
			name:  "Should not return error if the validation succeed",
			mocks: mocks{},
			args: args{
				values: map[string]interface{}{"name": "John", "age": 30},
				schemaJSON: []byte(`
				{
					"$schema": "http://json-schema.org/draft-07/schema#",
					"$id": "http://example.com/person.schema.json",
					"type": "object",
					"properties": {
						"name": {
							"type": "string"
						},
							"age": {
								"type": "integer",
								"minimum": 0
							}
						},
						"required": ["name"]
						}`,
				),
			},
			wantErr: "",
		},
		{
			name:  "Should describe the schema violation",
			mocks: mocks{},
			args: args{
				values: map[string]interface{}{"name": "John", "age": -1},
				schemaJSON: []byte(`
					{
						"$schema": "http://json-schema.org/draft-07/schema#",
						"$id": "http://example.com/person.schema.json",
						"type": "object",
						"properties": {
							"name": {
								"type": "string"
							},
							"age": {
								"type": "integer",
								"minimum": 0
							}
						},
						"required": ["name"]
					}`,
				),
			},
			wantErr: "- age: Must be greater than or equal to 0\n",
		},
		{
			name:  "Should propagate error if marshalling throws one",
			mocks: mocks{},
			args: args{
				values: map[string]interface{}{
					"name": make(chan int),
				},
				schemaJSON: []byte(`
					{
						"$schema": "http://json-schema.org/draft-07/schema#",
						"$id": "http://example.com/person.schema.json",
						"type": "object",
						"properties": {
							"name": {
								"type": "string"
							},
							"age": {
								"type": "integer",
								"minimum": 0
							}
						},
						"required": ["name"]
					}`,
				),
			},
			wantErr: "error marshaling into JSON: json: unsupported type: chan int",
		},
		{
			name:  "Should pass validation with empty schema and nil values",
			mocks: mocks{},
			args: args{
				values:     nil,
				schemaJSON: []byte(`{}`),
			},
			wantErr: "",
		},
		{
			name: "Should propagate error when YAML to JSON conversion throws it",
			mocks: mocks{
				yamlToJsonErr: errors.New("some-yaml-to-json-error"),
			},
			args: args{
				values:     nil,
				schemaJSON: []byte(`{}`),
			},
			wantErr: "some-yaml-to-json-error",
		},
		{
			name: "Should propagate error when json validation throws one",
			mocks: mocks{
				jsonSchemaValidationErr: errors.New("some-json-validation-error"),
			},
			args: args{
				values:     nil,
				schemaJSON: []byte(`{}`),
			},
			wantErr: "some-json-validation-error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mocks.yamlToJsonErr != nil {
				mockYamlToJson(t, tt.mocks.yamlToJsonErr)
			}
			if tt.mocks.jsonSchemaValidationErr != nil {
				mockValidateJson(t, tt.mocks.jsonSchemaValidationErr)
			}

			err := Validate(tt.args.values, tt.args.schemaJSON)

			if len(tt.wantErr) == 0 {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}

func mockYamlToJson(t *testing.T, err error) {
	originalValue := yamlToJson
	yamlToJson = func(y []byte) ([]byte, error) {
		return nil, err
	}
	t.Cleanup(func() { yamlToJson = originalValue })
}

func mockValidateJson(t *testing.T, err error) {
	originalValue := validateJson
	validateJson = func(ls, ld gojsonschema.JSONLoader) (*gojsonschema.Result, error) {
		return nil, err
	}
	t.Cleanup(func() { validateJson = originalValue })
}
