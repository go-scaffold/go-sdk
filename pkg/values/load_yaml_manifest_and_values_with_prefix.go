package values

import (
	"github.com/pasdam/go-template-map-loader/pkg/tm"
)

func LoadYamlManifestAndValuesWithPrefix(manifestPrefix string, manifestPath string, valuesPrefix string, valuesPaths ...string) (map[string]interface{}, error) {
	manifest, err := LoadYamlFilesWithPrefix(manifestPrefix, manifestPath)
	if err != nil {
		return nil, err
	}

	values, err := LoadYamlFilesWithPrefix(valuesPrefix, valuesPaths...)
	if err != nil {
		return nil, err
	}

	return tm.MergeMaps(manifest, values), nil
}
