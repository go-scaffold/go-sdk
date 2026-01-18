package templates

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func Test_ProcessTemplate_fail_shouldPropagateErrorIfReaderThrowsIt(t *testing.T) {
	file, err := os.Open("testdata/template_file.tpl")
	assert.Nil(t, err)
	defer file.Close()
	funcMap := template.FuncMap{}
	expectedErr := errors.New("some-read-error")
	mockReadAll(t, expectedErr)

	reader, err := ProcessTemplate(file, "invalid-data", funcMap)

	assert.Equal(t, expectedErr, err)
	assert.Nil(t, reader)
}

func Test_ProcessTemplate_fail_shouldReturnErrorIfApplyingTheTemplateFailed(t *testing.T) {
	file, err := os.Open("testdata/template_file.tpl")
	assert.Nil(t, err)
	defer file.Close()
	funcMap := template.FuncMap{}

	reader, err := ProcessTemplate(file, "invalid-data", funcMap)

	assert.NotNil(t, err)
	assert.Nil(t, reader)
}

func Test_ProcessTemplate_success_shouldCreateAReaderForTheGeneratedContent(t *testing.T) {
	file, err := os.Open("testdata/template_file.tpl")
	assert.Nil(t, err)
	defer file.Close()
	funcMap := template.FuncMap{
		"Bold": func(value string) string {
			return fmt.Sprintf("*%s*", value)
		},
	}

	reader, err := ProcessTemplate(file, struct{ Text string }{Text: "test"}, funcMap)

	assert.Nil(t, err)
	readContent, err := io.ReadAll(reader)
	assert.Nil(t, err)
	assert.Equal(t, "This is a *test*\n", string(readContent))
}

func mockReadAll(t *testing.T, err error) {
	originalValue := readAll
	readAll = func(r io.Reader) ([]byte, error) {
		return nil, err
	}
	t.Cleanup(func() { readAll = originalValue })
}

func Test_ProcessTemplateWithBaseTemplate_success_shouldUseCommonTemplatesFromBase(t *testing.T) {
	funcMap := template.FuncMap{}

	// Create a base template with common templates
	baseTemplate := template.New("base").Funcs(funcMap)
	baseTemplate, err := baseTemplate.New("header").Parse("[HEADER]")
	assert.NoError(t, err)
	baseTemplate, err = baseTemplate.New("footer").Parse("[FOOTER]")
	assert.NoError(t, err)

	templateContent := strings.NewReader("{{ template \"header\" }} Content {{ template \"footer\" }}")
	data := struct{}{}

	reader, err := ProcessTemplateWithBaseTemplate(templateContent, data, funcMap, nil, baseTemplate)

	assert.Nil(t, err)
	readContent, err := io.ReadAll(reader)
	assert.Nil(t, err)
	assert.Equal(t, "[HEADER] Content [FOOTER]", string(readContent))
}

func Test_ProcessTemplateWithBaseTemplate_success_shouldWorkWithNilBaseTemplate(t *testing.T) {
	funcMap := template.FuncMap{
		"Upper": strings.ToUpper,
	}

	templateContent := strings.NewReader("Hello {{ Upper .Name }}")
	data := struct{ Name string }{Name: "world"}

	reader, err := ProcessTemplateWithBaseTemplate(templateContent, data, funcMap, nil, nil)

	assert.Nil(t, err)
	readContent, err := io.ReadAll(reader)
	assert.Nil(t, err)
	assert.Equal(t, "Hello WORLD", string(readContent))
}

func Test_ProcessTemplateWithBaseTemplate_fail_shouldPropagateErrorIfReaderThrowsIt(t *testing.T) {
	funcMap := template.FuncMap{}
	baseTemplate := template.New("base").Funcs(funcMap)
	expectedErr := errors.New("some-read-error")
	mockReadAll(t, expectedErr)

	reader, err := ProcessTemplateWithBaseTemplate(strings.NewReader("test"), struct{}{}, funcMap, nil, baseTemplate)

	assert.Equal(t, expectedErr, err)
	assert.Nil(t, reader)
}

func Test_ProcessTemplateWithBaseTemplate_fail_shouldReturnErrorIfTemplateReferencesNonExistent(t *testing.T) {
	funcMap := template.FuncMap{}
	baseTemplate := template.New("base").Funcs(funcMap)

	templateContent := strings.NewReader("{{ template \"nonexistent\" }}")
	data := struct{}{}

	reader, err := ProcessTemplateWithBaseTemplate(templateContent, data, funcMap, nil, baseTemplate)

	assert.NotNil(t, err)
	assert.Nil(t, reader)
}

func Test_ProcessTemplateWithBaseTemplate_success_shouldWorkWithTemplateAwareFunctions(t *testing.T) {
	funcMap := template.FuncMap{}
	templateAwareFnGen := TemplateAwareFuncMap{
		"RenderTemplate": func(t *template.Template) any {
			return func(name string) string {
				var buf strings.Builder
				t.ExecuteTemplate(&buf, name, nil)
				return buf.String()
			}
		},
	}

	// Create a base template with a common template
	baseTemplate := template.New("base").Funcs(funcMap)
	baseTemplate, err := baseTemplate.New("snippet").Parse("SNIPPET_CONTENT")
	assert.NoError(t, err)

	templateContent := strings.NewReader("Result: {{ RenderTemplate \"snippet\" }}")
	data := struct{}{}

	reader, err := ProcessTemplateWithBaseTemplate(templateContent, data, funcMap, templateAwareFnGen, baseTemplate)

	assert.Nil(t, err)
	readContent, err := io.ReadAll(reader)
	assert.Nil(t, err)
	assert.Equal(t, "Result: SNIPPET_CONTENT", string(readContent))
}
