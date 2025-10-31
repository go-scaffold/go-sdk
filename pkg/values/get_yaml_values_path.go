package values

func GetYamlValuesPath(dir string) (string, error) {
	return GetYamlPath(dir, ValuesPrefix)
}
