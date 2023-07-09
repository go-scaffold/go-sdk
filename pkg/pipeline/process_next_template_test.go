package pipeline

import (
	"errors"
	"go-sdk/pkg/testhelpers"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/pasdam/go-io-utilx/pkg/ioutilx"
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
			templateProcessor := &templateProcessorMock{}
			data := make(map[string]string)
			funcMap := make(template.FuncMap)

			var nextTemplate *templateMock
			if tt.mocks.nextTemplateErr == nil {
				nextTemplate = &templateMock{}
				templateReader := io.NopCloser(strings.NewReader(tt.mocks.templateContent))
				nextTemplate.On("Reader").Return(templateReader)
				nextTemplate.On("Path").Return(tt.wantPath)
				templateProcessor.On("NextTemplate").Return(nextTemplate, nil)
				mockProcessTemplate(t, templateReader, data, funcMap, tt.wantContent, tt.mocks.renderTemplateErr)

			} else {
				templateProcessor.On("NextTemplate").Return(nil, tt.mocks.nextTemplateErr)
			}

			got, err := processNextTemplate(templateProcessor, data, funcMap)

			if tt.wantErr == nil {
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantPath, got.Path)
				assert.Equal(t, tt.wantContent, ioutilx.ReaderToString(got.Reader))

			} else {
				assert.Nil(t, got)
			}
			testhelpers.AssertEqualErrors(t, tt.wantErr, err)
		})
	}
}

func mockProcessTemplate(t *testing.T, expectedReader io.Reader, expectedData interface{}, expectedFuncMap template.FuncMap, content string, err error) {
	originalValue := _processTemplate
	_processTemplate = func(gotReader io.Reader, gotData interface{}, gotFuncMap template.FuncMap) (io.Reader, error) {
		assert.Equal(t, expectedReader, gotReader)
		assert.Equal(t, expectedData, gotData)
		assert.Equal(t, expectedFuncMap, gotFuncMap)
		if err == nil {
			return strings.NewReader(content), nil
		}
		return nil, err
	}
	t.Cleanup(func() { _processTemplate = originalValue })
}
