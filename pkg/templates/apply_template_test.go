package templates

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func Test_applyTemplate_Fail_ShouldReturnErrorIfItFailsToExecuteTheTemplate(t *testing.T) {
	funcMap := template.FuncMap{}
	result, err := applyTemplate("This is a {{ .NotExistingProperty }}", "invalid_config", funcMap, nil)

	assert.NotNil(t, err)
	assert.Empty(t, result)
}

func Test_applyTemplate_Fail_ShouldReturnErrorIfTemplateIsInvalid(t *testing.T) {
	funcMap := template.FuncMap{}
	data := struct{ CustomProperty string }{CustomProperty: "*test*"}
	result, err := applyTemplate("This is a {{ .CustomProperty } with invalid template", data, funcMap, nil)

	assert.NotNil(t, err)
	assert.Empty(t, result)
}

func Test_applyTemplate_Success_ShouldCorrectlyGenerateOutputContentFromTemplate(t *testing.T) {
	funcMap := template.FuncMap{
		"Bold": func(value string) string {
			return fmt.Sprintf("*%s*", value)
		},
	}
	templateAwareFuncMap := TemplateAwareFuncMap{
		"TemplateAware": func(t *template.Template) any {
			return func(name string, value string) string {
				var result bytes.Buffer
				t.ExecuteTemplate(&result, name, value)
				return result.String()
			}
		},
	}
	data := struct{ CustomProperty string }{CustomProperty: "test"}
	result, err := applyTemplate("{{ define \"tt\" }}--{{ . }}--{{ end }}This is a {{ Bold .CustomProperty }} {{ TemplateAware \"tt\" .CustomProperty }}", data, funcMap, templateAwareFuncMap)

	assert.Nil(t, err)
	actualContent, err := io.ReadAll(result)
	assert.NoError(t, err)
	assert.Equal(t, "This is a *test* --test--", string(actualContent))
}

func Test_applyTemplateWithBase_Success_ShouldUseCommonTemplatesFromBase(t *testing.T) {
	funcMap := template.FuncMap{
		"Bold": func(value string) string {
			return fmt.Sprintf("*%s*", value)
		},
	}

	// Create a base template with a common template defined
	baseTemplate := template.New("base").Funcs(funcMap)
	baseTemplate, err := baseTemplate.New("header").Parse("=== {{ . }} ===")
	assert.NoError(t, err)
	baseTemplate, err = baseTemplate.New("footer").Parse("--- {{ . }} ---")
	assert.NoError(t, err)

	data := struct{ Title string }{Title: "Hello"}
	result, err := applyTemplateWithBase("{{ template \"header\" .Title }} Content {{ template \"footer\" .Title }}", data, funcMap, nil, baseTemplate)

	assert.Nil(t, err)
	actualContent, err := io.ReadAll(result)
	assert.NoError(t, err)
	assert.Equal(t, "=== Hello === Content --- Hello ---", string(actualContent))
}

func Test_applyTemplateWithBase_Success_ShouldWorkWithNilBaseTemplate(t *testing.T) {
	funcMap := template.FuncMap{}
	data := struct{ Text string }{Text: "test"}
	result, err := applyTemplateWithBase("Hello {{ .Text }}", data, funcMap, nil, nil)

	assert.Nil(t, err)
	actualContent, err := io.ReadAll(result)
	assert.NoError(t, err)
	assert.Equal(t, "Hello test", string(actualContent))
}

func Test_applyTemplateWithBase_Success_ShouldCloneBaseTemplateAndNotModifyOriginal(t *testing.T) {
	funcMap := template.FuncMap{}

	// Create a base template
	baseTemplate := template.New("base").Funcs(funcMap)
	baseTemplate, err := baseTemplate.New("common").Parse("COMMON")
	assert.NoError(t, err)

	// Get initial template count
	initialTemplates := baseTemplate.Templates()

	data := struct{}{}
	// This will parse new content into a cloned template
	result, err := applyTemplateWithBase("{{ define \"new\" }}NEW{{ end }}{{ template \"common\" }}", data, funcMap, nil, baseTemplate)

	assert.Nil(t, err)
	actualContent, err := io.ReadAll(result)
	assert.NoError(t, err)
	assert.Equal(t, "COMMON", string(actualContent))

	// Verify base template was not modified (still has same number of templates)
	assert.Equal(t, len(initialTemplates), len(baseTemplate.Templates()))
}

func Test_applyTemplateWithBase_Fail_ShouldReturnErrorIfBaseCloneFails(t *testing.T) {
	funcMap := template.FuncMap{}

	// Create a base template with a common template
	baseTemplate := template.New("base").Funcs(funcMap)
	baseTemplate, err := baseTemplate.New("common").Parse("COMMON")
	assert.NoError(t, err)

	data := struct{}{}
	// Reference a non-existent template to cause an execution error
	result, err := applyTemplateWithBase("{{ template \"nonexistent\" }}", data, funcMap, nil, baseTemplate)

	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func Test_applyTemplateWithBase_Success_ShouldWorkWithTemplateAwareFunctions(t *testing.T) {
	funcMap := template.FuncMap{}
	templateAwareFuncMap := TemplateAwareFuncMap{
		"CallTemplate": func(t *template.Template) any {
			return func(name string, data interface{}) string {
				var result bytes.Buffer
				t.ExecuteTemplate(&result, name, data)
				return result.String()
			}
		},
	}

	// Create a base template with a common template
	baseTemplate := template.New("base").Funcs(funcMap)
	baseTemplate, err := baseTemplate.New("greet").Parse("Hello, {{ . }}!")
	assert.NoError(t, err)

	data := struct{ Name string }{Name: "World"}
	result, err := applyTemplateWithBase("{{ CallTemplate \"greet\" .Name }}", data, funcMap, templateAwareFuncMap, baseTemplate)

	assert.Nil(t, err)
	actualContent, err := io.ReadAll(result)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", string(actualContent))
}
