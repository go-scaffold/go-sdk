package values

import (
	"github.com/pasdam/go-template-map-loader/pkg/tm"
)

func LoadYamlFilesWithPrefix(prefix string, paths ...string) (map[string]interface{}, error) {
	valuesData := make([]map[string]interface{}, 0, len(paths))

	for _, path := range paths {
		data, err := tm.LoadYamlFile(path)
		if err != nil {
			return nil, err
		}
		valuesData = append(valuesData, data)
	}

	return tm.WithPrefix(prefix, tm.MergeMaps(valuesData...)), nil
}
