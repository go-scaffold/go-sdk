package pipeline

import (
	"errors"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/pasdam/go-template-map-loader/pkg/tm"
	"github.com/pasdam/go-test-utils/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_pipeline_Process(t *testing.T) {
	type fields struct {
		prefixData           string
		prefixMetadata       string
		withDataPreprocessor bool
	}
	type mocks struct {
		nextTemplateRes       []*nextTemplateResult
		collectingErrs        []error
		dataPreprocessorError error
	}
	tests := []struct {
		name    string
		fields  fields
		mocks   mocks
		wantErr error
	}{
		{
			name:    "Should not return error if next template returns EoF",
			mocks:   mocks{},
			wantErr: nil,
		},
		{
			name: "Should not return error if next template returns valid values",
			mocks: mocks{
				nextTemplateRes: []*nextTemplateResult{
					{
						data: &Template{
							Reader: io.NopCloser(strings.NewReader("")),
						},
						err: nil,
					},
					{
						data: &Template{
							Reader: io.NopCloser(strings.NewReader("")),
						},
						err: nil,
					},
					{
						data: nil,
						err:  io.EOF,
					},
				},
				collectingErrs: []error{
					nil,
					nil,
					nil,
				},
			},
			wantErr: nil,
		},
		{
			name: "Should not return error if next template returns valid values and pipeline is using prefixes and a data preprocessor",
			fields: fields{
				prefixData:           "some-data-prefix",
				prefixMetadata:       "some-metadata-prefix",
				withDataPreprocessor: true,
			},
			mocks: mocks{
				nextTemplateRes: []*nextTemplateResult{
					{
						data: &Template{
							Reader: io.NopCloser(strings.NewReader("")),
						},
						err: nil,
					},
					{
						data: &Template{
							Reader: io.NopCloser(strings.NewReader("")),
						},
						err: nil,
					},
					{
						data: nil,
						err:  io.EOF,
					},
				},
				collectingErrs: []error{
					nil,
					nil,
					nil,
				},
			},
			wantErr: nil,
		},
		{
			name: "Should propagate error if next template returns one",
			mocks: mocks{
				nextTemplateRes: []*nextTemplateResult{
					{
						data: &Template{
							Reader: io.NopCloser(strings.NewReader("")),
						},
						err: nil,
					},
					{
						data: &Template{
							Reader: io.NopCloser(strings.NewReader("")),
						},
						err: nil,
					},
					{
						data: nil,
						err:  errors.New("some-unexpected-error"),
					},
				},
				collectingErrs: []error{
					nil,
					nil,
					nil,
				},
			},
			wantErr: errors.New("some-unexpected-error"),
		},
		{
			name: "Should propagate error if data processor returns one",
			fields: fields{
				withDataPreprocessor: true,
			},
			mocks: mocks{
				dataPreprocessorError: errors.New("some-data-processor-error"),
			},
			wantErr: errors.New("some-data-processor-error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templateProvider := &templateProviderMock{}
			data := map[string]interface{}{}
			metadata := map[string]interface{}{}
			expectedData := tm.MergeMaps(
				tm.WithPrefix(tt.fields.prefixData, data),
				tm.WithPrefix(tt.fields.prefixMetadata, metadata),
			)
			functions := make(template.FuncMap)
			collector := &collectorMock{}
			p := &pipeline{
				collector:        collector,
				functions:        functions,
				prefixData:       tt.fields.prefixData,
				prefixMetadata:   tt.fields.prefixMetadata,
				templateProvider: templateProvider,
			}
			if tt.fields.withDataPreprocessor {
				oldExpectedData := expectedData
				expectedDataWithPrefix := tm.WithPrefix("some-preprocessor-prefix", expectedData)
				p.dataPreprocessor = func(m map[string]interface{}) (map[string]interface{}, error) {
					assert.Equal(t, oldExpectedData, m)
					if tt.mocks.dataPreprocessorError != nil {
						return nil, tt.mocks.dataPreprocessorError
					}
					return expectedDataWithPrefix, nil
				}
				expectedData = expectedDataWithPrefix
			}
			mockProcessNextTemplate(t, templateProvider, expectedData, functions, tt.mocks.nextTemplateRes)
			assert.Equal(t, len(tt.mocks.nextTemplateRes), len(tt.mocks.collectingErrs))
			for i := 0; i < len(tt.mocks.nextTemplateRes); i++ {
				if tt.mocks.nextTemplateRes[i].err == nil {
					collector.On("Collect", tt.mocks.nextTemplateRes[i].data).Return(tt.mocks.collectingErrs[i])
				}
			}

			err := p.Process(metadata, data)

			testutils.AssertEqualErrors(t, tt.wantErr, err)
		})
	}
}

type nextTemplateResult struct {
	data *Template
	err  error
}

func mockProcessNextTemplate(t *testing.T, expectedProcessor TemplateProvider, expectedData interface{}, expectedFuncMap template.FuncMap, nextTemplateRes []*nextTemplateResult) {
	originalValue := _processNextTemplate
	count := 0
	_processNextTemplate = func(gotProcessor TemplateProvider, gotData interface{}, gotFuncMap template.FuncMap) (*Template, error) {
		assert.Equal(t, expectedProcessor, gotProcessor)
		assert.Equal(t, expectedData, gotData)
		assert.Equal(t, expectedFuncMap, gotFuncMap)

		if len(nextTemplateRes) == 0 {
			return nil, io.EOF
		}
		assert.True(t, count < len(nextTemplateRes))
		res := nextTemplateRes[count]
		count++
		return res.data, res.err
	}
	t.Cleanup(func() { _processNextTemplate = originalValue })
}
