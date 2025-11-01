package values

import (
	"fmt"

	"github.com/pasdam/go-template-map-loader/pkg/tm"
)

const (
	defaultManifestPrefix   = "Manifest"
	defaultManifestBasename = "Manifest"
	defaultValuesPrefix     = "Values"
	defaultValuesBasename   = "values"
)

// Loader represents a configuration for loading YAML files with customizable prefixes and basenames
type Loader struct {
	manifestPrefix   string
	manifestBasename string
	valuesPrefix     string
	valuesBasename   string
}

// NewLoader creates a new Loader with default values:
// manifestPrefix = "Manifest"
// manifestBasename = "Manifest"
// valuesPrefix = "Values"
// valuesBasename = "Values"
func NewLoader() *Loader {
	return &Loader{
		manifestPrefix:   defaultManifestPrefix,
		manifestBasename: defaultManifestBasename,
		valuesPrefix:     defaultValuesPrefix,
		valuesBasename:   defaultValuesBasename,
	}
}

// NewLoaderWithValues creates a new Loader with custom values for all fields
func NewLoaderWithValues(manifestPrefix, manifestBasename, valuesPrefix, valuesBasename string) *Loader {
	return &Loader{
		manifestPrefix:   manifestPrefix,
		manifestBasename: manifestBasename,
		valuesPrefix:     valuesPrefix,
		valuesBasename:   valuesBasename,
	}
}

func (l *Loader) LoadYAMLs(manifestDir string, additionalValueFiles []string) (map[string]interface{}, error) {
	manifestPath, err := GetYamlPath(manifestDir, l.manifestBasename)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while getting the manifest path: %s", err.Error())
	}
	manifest, err := LoadYamlFilesWithPrefix(l.manifestPrefix, manifestPath)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while reading the manifest file: %s", err.Error())
	}

	valuesPath, err := GetYamlPath(manifestDir, l.valuesBasename)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while getting the value path: %s", err.Error())
	}

	valuesPaths := make([]string, 0, len(additionalValueFiles)+1)
	valuesPaths = append(valuesPaths, valuesPath)
	valuesPaths = append(valuesPaths, additionalValueFiles...)
	data, err := LoadYamlFilesWithPrefix(l.valuesPrefix, valuesPaths...)
	if err != nil {
		return nil, fmt.Errorf("error while loading data: %s", err.Error())
	}

	return tm.MergeMaps(manifest, data), nil
}
