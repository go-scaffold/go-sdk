package values

func LoadYamlValues(paths []string) (map[string]interface{}, error) {
	return LoadYamlFilesWithPrefix(defaultValuesPrefix, paths...)
}
