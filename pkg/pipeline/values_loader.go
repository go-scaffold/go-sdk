package pipeline

type ValuesLoader interface {
	LoadYamlFilesWithPrefix(prefix string, paths ...string) (map[string]interface{}, error)
}
