package pipeline

import (
	"errors"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/pasdam/go-test-utils/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_pipeline_Process(t *testing.T) {
	type mocks struct {
		nextTemplateRes    []*nextTemplateResult
		postProcessingErrs []error
	}
	tests := []struct {
		name    string
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
				postProcessingErrs: []error{
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
				postProcessingErrs: []error{
					nil,
					nil,
					nil,
				},
			},
			wantErr: errors.New("some-unexpected-error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templateProvider := &templateProviderMock{}
			data := map[string]interface{}{}
			functions := make(template.FuncMap)
			postProcessor := &postProcessorMock{}
			p := &pipeline{
				data:             data,
				functions:        functions,
				postProcessor:    postProcessor,
				templateProvider: templateProvider,
			}
			mockProcessNextTemplate(t, templateProvider, data, functions, tt.mocks.nextTemplateRes)
			assert.Equal(t, len(tt.mocks.nextTemplateRes), len(tt.mocks.postProcessingErrs))
			for i := 0; i < len(tt.mocks.nextTemplateRes); i++ {
				if tt.mocks.nextTemplateRes[i].err == nil {
					postProcessor.On("Process", tt.mocks.nextTemplateRes[i].data).Return(tt.mocks.postProcessingErrs[i])
				}
			}

			err := p.Process()

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
