package values

func GetYamlManifestPath(dir string) (string, error) {
	return GetYamlPath(dir, ManifestPrefix)
}
