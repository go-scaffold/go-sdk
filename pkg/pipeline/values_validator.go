package pipeline

type ValuesValidator interface {
	ValidateYaml(values map[string]interface{}, schemaJSON []byte) error
}
