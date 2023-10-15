package pipeline

type DataPreprocessor func(map[string]interface{}) (map[string]interface{}, error)
