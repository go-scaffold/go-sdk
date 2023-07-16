package values

func LoadYamlManifestAndValues(manifestPath string, valuesPaths ...string) (map[string]interface{}, error) {
	return LoadYamlManifestAndValuesWithPrefix(defaultManifestPrefix, manifestPath, defaultValuesPrefix, valuesPaths...)
}
