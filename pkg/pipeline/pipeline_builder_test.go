package pipeline

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func Test_builder_Build(t *testing.T) {
	funcMap := make(template.FuncMap)
	funcMap["test"] = Test_builder_Build
	type fields struct {
		dataPrefix        string
		data              map[string]interface{}
		metadataPrefix    string
		metadata          map[string]interface{}
		functions         template.FuncMap
		templateProcessor TemplateProcessor
		postProcessors    []PostProcessor
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr string
	}{
		{
			name: "Should return error if data and metadata are nill",
			fields: fields{
				data:              nil,
				metadata:          nil,
				functions:         funcMap,
				templateProcessor: &templateProcessorMock{},
				postProcessors: []PostProcessor{
					&postProcessorMock{},
				},
			},
			wantErr: "no data specified in the context",
		},
		{
			name: "Should return error if functions map is nill",
			fields: fields{
				data: map[string]interface{}{
					"kd1": "vd1",
				},
				metadata: map[string]interface{}{
					"km2": "vm2",
				},
				functions:         nil,
				templateProcessor: &templateProcessorMock{},
				postProcessors: []PostProcessor{
					&postProcessorMock{},
				},
			},
			wantErr: "no functions specified in the context",
		},
		{
			name: "Should return error if functions map is empty",
			fields: fields{
				data: map[string]interface{}{
					"kd1": "vd1",
				},
				metadata: map[string]interface{}{
					"km2": "vm2",
				},
				functions:         make(template.FuncMap),
				templateProcessor: &templateProcessorMock{},
				postProcessors: []PostProcessor{
					&postProcessorMock{},
				},
			},
			wantErr: "no functions specified in the context",
		},
		{
			name: "Should return error if template processor is nil",
			fields: fields{
				data: map[string]interface{}{
					"kd1": "vd1",
				},
				metadata: map[string]interface{}{
					"km2": "vm2",
				},
				functions:         funcMap,
				templateProcessor: nil,
				postProcessors: []PostProcessor{
					&postProcessorMock{},
				},
			},
			wantErr: "no template processor specified for the pipeline",
		},
		{
			name: "Should return error if there are no post processors",
			fields: fields{
				data: map[string]interface{}{
					"kd1": "vd1",
				},
				metadata: map[string]interface{}{
					"km2": "vm2",
				},
				functions:         funcMap,
				templateProcessor: &templateProcessorMock{},
				postProcessors:    nil,
			},
			wantErr: "no post processor specified for the pipeline",
		},
		{
			name: "Should create pipeline with one post processor",
			fields: fields{
				data: map[string]interface{}{
					"kd1": "vd1",
					"kd2": "vd2",
				},
				metadata: map[string]interface{}{
					"km1": "vm1",
					"km2": "vm2",
				},
				functions:         funcMap,
				templateProcessor: &templateProcessorMock{},
				postProcessors: []PostProcessor{
					&postProcessorMock{},
				},
			},
		},
		{
			name: "Should create pipeline with one post processor and prefixes for data and metadata",
			fields: fields{
				dataPrefix: "CustomData",
				data: map[string]interface{}{
					"kd1": "vd1",
					"kd2": "vd2",
				},
				metadataPrefix: "CustomMetadata",
				metadata: map[string]interface{}{
					"km1": "vm1",
					"km2": "vm2",
				},
				functions:         funcMap,
				templateProcessor: &templateProcessorMock{},
				postProcessors: []PostProcessor{
					&postProcessorMock{},
				},
			},
		},
		{
			name: "Should create pipeline with 2 post processors",
			fields: fields{
				data: map[string]interface{}{
					"kd3": "vd3",
					"kd4": "vd4",
				},
				metadata: map[string]interface{}{
					"km3": "vm3",
					"km4": "vm4",
				},
				functions:         funcMap,
				templateProcessor: &templateProcessorMock{},
				postProcessors: []PostProcessor{
					&postProcessorMock{},
					&postProcessorMock{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mBuilder := NewBuilder()
			var dBuilder DataBuilder
			if tt.fields.metadataPrefix == "" {
				dBuilder = mBuilder.WithMetadata(tt.fields.metadata)
			} else {
				dBuilder = mBuilder.WithMetadataWithPrefix(tt.fields.metadataPrefix, tt.fields.metadata)
			}
			var fBuilder FunctionsBuilder
			if tt.fields.dataPrefix == "" {
				fBuilder = dBuilder.WithData(tt.fields.data)
			} else {
				fBuilder = dBuilder.WithDataWithPrefix(tt.fields.dataPrefix, tt.fields.data)
			}

			builder := fBuilder.
				WithFunctions(tt.fields.functions).
				WithTemplateProcessor(tt.fields.templateProcessor)

			for _, processor := range tt.fields.postProcessors {
				builder = builder.AddResultProcessor(processor)
			}

			got, err := builder.Build()

			if len(tt.wantErr) == 0 {
				assert.NoError(t, err)
				assert.NotNil(t, got)

				gotP := got.(*pipeline)
				assert.Len(t, gotP.data, 2)
				if tt.fields.dataPrefix == "" {
					assert.Equal(t, tt.fields.data, gotP.data["Values"])
				} else {
					assert.Equal(t, tt.fields.data, gotP.data[tt.fields.dataPrefix])
				}
				if tt.fields.metadataPrefix == "" {
					assert.Equal(t, tt.fields.metadata, gotP.data["Manifest"])
				} else {
					assert.Equal(t, tt.fields.metadata, gotP.data[tt.fields.metadataPrefix])
				}
				assert.Equal(t, tt.fields.functions, gotP.functions)
				assert.Equal(t, tt.fields.templateProcessor, gotP.templateProcessor)

				postProcessorStep := gotP.postProcessingSteps
				assert.Equal(t, tt.fields.postProcessors[0], postProcessorStep.processor)
				for i := 1; i < len(tt.fields.postProcessors); i++ {
					assert.NotNil(t, postProcessorStep.nextStep)

					postProcessorStep = postProcessorStep.nextStep
					assert.Equal(t, tt.fields.postProcessors[i], postProcessorStep.processor)
				}
				assert.Nil(t, postProcessorStep.nextStep)

			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr)
				assert.Nil(t, got)
			}
		})
	}
}
