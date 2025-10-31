package values

import (
	"os"
	"path/filepath"
)

// GetYamlPath returns the full path of a YAML file given a directory path and
// a yaml basename, regardless if it ends with .yaml or .yml.
// It checks for both extensions in the specified directory and returns the
// first one that exists, or an empty string if none is found.
func GetYamlPath(dirPath, baseName string) (string, error) {
	// Check for .yaml extension first
	yamlPath := filepath.Join(dirPath, baseName+".yaml")
	if _, err := os.Stat(yamlPath); err == nil {
		return yamlPath, nil
	} else if !os.IsNotExist(err) {
		// If there's an error other than "not exist", return it
		return "", err
	}

	// Then check for .yml extension
	ymlPath := filepath.Join(dirPath, baseName+".yml")
	if _, err := os.Stat(ymlPath); err == nil {
		return ymlPath, nil
	} else if !os.IsNotExist(err) {
		// If there's an error other than "not exist", return it
		return "", err
	}

	// Neither extension exists
	return "", nil
}
