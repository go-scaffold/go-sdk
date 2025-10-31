package values

func LoadYamlManifestAndValues(manifestPath string, valuesPaths ...string) (map[string]interface{}, error) {
	return LoadYamlManifestAndValuesWithPrefix(ManifestPrefix, manifestPath, ValuesPrefix, valuesPaths...)
}
