package pipeline

type MetadataBuilder interface {
	WithMetadata(data map[string]interface{}) DataBuilder
	WithMetadataWithPrefix(prefix string, data map[string]interface{}) DataBuilder
}
