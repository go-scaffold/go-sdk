package pipeline

import (
	"errors"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/go-scaffold/go-sdk/v2/pkg/templates"
	"github.com/pasdam/go-io-utilx/pkg/ioutilx"
	"github.com/pasdam/go-utils/pkg/assertutils"
	"github.com/stretchr/testify/assert"
)

func Test_processNextTemplate(t *testing.T) {
	type mocks struct {
		templateContent   string
		nextTemplateErr   error
		renderTemplateErr error
	}
	tests := []struct {
		name        string
		mocks       mocks
		wantContent string
		wantPath    string
		wantErr     error
	}{
		{
			name:        "Should return processed data if no error occurs",
			mocks:       mocks{},
			wantPath:    "some-path",
			wantContent: "some-reader-content",
		},
		{
			name: "Should propagate the error if next template returns on",
			mocks: mocks{
				nextTemplateErr: errors.New("some next template error"),
			},
			wantErr: errors.New("some next template error"),
		},
		{
			name: "Should propagate the error if render template returns on",
			mocks: mocks{
				renderTemplateErr: errors.New("some render template error"),
			},
			wantErr: errors.New("some render template error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templateProvider := &templateProviderMock{}
			data := make(map[string]string)
			funcMap := make(template.FuncMap)
			templateAwareFnGen := make(templates.TemplateAwareFuncMap)

			var nextTemplate *Template
			if tt.mocks.nextTemplateErr == nil {
				templateReader := io.NopCloser(strings.NewReader(tt.mocks.templateContent))
				nextTemplate = &Template{
					Reader: templateReader,
					Path:   tt.wantPath,
				}
				templateProvider.On("NextTemplate").Return(nextTemplate, nil)
				mockProcessTemplate(t, templateReader, data, funcMap, templateAwareFnGen, nil, tt.wantContent, tt.mocks.renderTemplateErr)

			} else {
				templateProvider.On("NextTemplate").Return(nil, tt.mocks.nextTemplateErr)
			}

			got, err := processNextTemplate(templateProvider, data, funcMap, templateAwareFnGen, nil)

			if tt.wantErr == nil {
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantPath, got.Path)
				assert.Equal(t, tt.wantContent, ioutilx.ReaderToString(got.Reader))

			} else {
				assert.Nil(t, got)
			}
			assertutils.AssertEqualErrors(t, tt.wantErr, err)
		})
	}
}

func mockProcessTemplate(t *testing.T, expectedReader io.Reader, expectedData interface{}, expectedFuncMap template.FuncMap, expectedTemplateAwareFnGen templates.TemplateAwareFuncMap, expectedBaseTemplate *template.Template, content string, err error) {
	originalValue := _processTemplate
	_processTemplate = func(gotReader io.Reader, gotData interface{}, gotFuncMap template.FuncMap, gotTemplateAwareFnGen templates.TemplateAwareFuncMap, gotBaseTemplate *template.Template) (io.Reader, error) {
		assert.Equal(t, expectedReader, gotReader)
		assert.Equal(t, expectedData, gotData)
		assert.Equal(t, expectedFuncMap, gotFuncMap)
		assert.Equal(t, expectedTemplateAwareFnGen, gotTemplateAwareFnGen)
		assert.Equal(t, expectedBaseTemplate, gotBaseTemplate)
		if err == nil {
			return strings.NewReader(content), nil
		}
		return nil, err
	}
	t.Cleanup(func() { _processTemplate = originalValue })
}

func Test_processNextTemplate_WithBaseTemplate(t *testing.T) {
	type mocks struct {
		templateContent   string
		nextTemplateErr   error
		renderTemplateErr error
	}
	tests := []struct {
		name         string
		mocks        mocks
		baseTemplate *template.Template
		wantContent  string
		wantPath     string
		wantErr      error
	}{
		{
			name:         "Should pass base template to process template",
			mocks:        mocks{},
			baseTemplate: template.Must(template.New("base").Parse("base content")),
			wantPath:     "some-path",
			wantContent:  "processed-with-base",
		},
		{
			name:         "Should work with nil base template",
			mocks:        mocks{},
			baseTemplate: nil,
			wantPath:     "another-path",
			wantContent:  "processed-without-base",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templateProvider := &templateProviderMock{}
			data := make(map[string]string)
			funcMap := make(template.FuncMap)
			templateAwareFnGen := make(templates.TemplateAwareFuncMap)

			var nextTemplate *Template
			if tt.mocks.nextTemplateErr == nil {
				templateReader := io.NopCloser(strings.NewReader(tt.mocks.templateContent))
				nextTemplate = &Template{
					Reader: templateReader,
					Path:   tt.wantPath,
				}
				templateProvider.On("NextTemplate").Return(nextTemplate, nil)
				mockProcessTemplate(t, templateReader, data, funcMap, templateAwareFnGen, tt.baseTemplate, tt.wantContent, tt.mocks.renderTemplateErr)

			} else {
				templateProvider.On("NextTemplate").Return(nil, tt.mocks.nextTemplateErr)
			}

			got, err := processNextTemplate(templateProvider, data, funcMap, templateAwareFnGen, tt.baseTemplate)

			if tt.wantErr == nil {
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantPath, got.Path)
				assert.Equal(t, tt.wantContent, ioutilx.ReaderToString(got.Reader))

			} else {
				assert.Nil(t, got)
			}
			assertutils.AssertEqualErrors(t, tt.wantErr, err)
		})
	}
}
