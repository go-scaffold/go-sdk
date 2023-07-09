package pipeline

type DataBuilder interface {
	WithData(data map[string]interface{}) FunctionsBuilder
	WithDataWithPrefix(prefix string, data map[string]interface{}) FunctionsBuilder
}
