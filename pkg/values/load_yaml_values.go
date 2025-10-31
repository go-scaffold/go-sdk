package values

func LoadYamlValues(paths []string) (map[string]interface{}, error) {
	return LoadYamlFilesWithPrefix(ValuesPrefix, paths...)
}
