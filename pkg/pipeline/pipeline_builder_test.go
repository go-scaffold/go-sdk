package pipeline

import (
	"errors"
	"testing"
	"text/template"

	"github.com/pasdam/go-utils/pkg/assertutils"
	"github.com/stretchr/testify/assert"
)

func Test_builder_Build(t *testing.T) {
	funcMap := make(template.FuncMap)
	funcMap["test"] = Test_builder_Build
	tests := []struct {
		name     string
		pipeline pipeline
		wantErr  error
	}{
		{
			name: "Should return error if functions map is nill",
			pipeline: pipeline{
				functions:        nil,
				templateProvider: &templateProviderMock{},
				collector:        &collectorMock{},
			},
			wantErr: errors.New("no functions specified in the context"),
		},
		{
			name: "Should return error if functions map is empty",
			pipeline: pipeline{
				functions:        make(template.FuncMap),
				templateProvider: &templateProviderMock{},
				collector:        &collectorMock{},
			},
			wantErr: errors.New("no functions specified in the context"),
		},
		{
			name: "Should return error if template processor is nil",
			pipeline: pipeline{
				functions:        funcMap,
				templateProvider: nil,
				collector:        &collectorMock{},
			},
			wantErr: errors.New("no template processor specified for the pipeline"),
		},
		{
			name: "Should return error if there is no collector",
			pipeline: pipeline{
				functions:        funcMap,
				templateProvider: &templateProviderMock{},
				collector:        nil,
			},
			wantErr: errors.New("no collector specified for the pipeline"),
		},
		{
			name: "Should create pipeline with a collector",
			pipeline: pipeline{
				functions:        funcMap,
				templateProvider: &templateProviderMock{},
				collector:        &collectorMock{},
			},
		},
		{
			name: "Should create pipeline with a collector and prefixes for data and metadata",
			pipeline: pipeline{
				functions:        funcMap,
				templateProvider: &templateProviderMock{},
				collector:        &collectorMock{},
			},
		},
		{
			name: "Should create pipeline with common templates provider",
			pipeline: pipeline{
				functions:              funcMap,
				templateProvider:       &templateProviderMock{},
				collector:              &collectorMock{},
				namedTemplatesProvider: &templateProviderMock{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewPipelineBuilder()

			builder = builder.
				WithDataPreprocessor(tt.pipeline.dataPreprocessor).
				WithFunctions(tt.pipeline.functions).
				WithTemplateProvider(tt.pipeline.templateProvider).
				WithCollector(tt.pipeline.collector).
				WithNamedTemplatesProvider(tt.pipeline.namedTemplatesProvider)

			expectedPipeline := tt.pipeline

			got, gotErr := builder.Build()

			assertutils.AssertEqualErrors(t, tt.wantErr, gotErr)
			if tt.wantErr == nil {
				assert.Equal(t, &expectedPipeline, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}
