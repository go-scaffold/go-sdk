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
