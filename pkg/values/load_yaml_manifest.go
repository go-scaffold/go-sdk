package values

func LoadYamlManifest(path string) (map[string]interface{}, error) {
	return LoadYamlFilesWithPrefix(ManifestPrefix, path)
}
