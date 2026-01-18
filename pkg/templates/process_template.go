package templates

import (
	"io"
	"log"
	"text/template"
)

var readAll = io.ReadAll

// ProcessTemplate processes the template using the specified data
func ProcessTemplate(reader io.Reader, data interface{}, funcMap template.FuncMap) (io.Reader, error) {
	return ProcessTemplateWithTemplateAware(reader, data, funcMap, nil)
}

// ProcessTemplateWithTemplateAware processes the template using the specified data
func ProcessTemplateWithTemplateAware(reader io.Reader, data interface{}, funcMap template.FuncMap, templateAwareFuncGenerators TemplateAwareFuncMap) (io.Reader, error) {
	return ProcessTemplateWithBaseTemplate(reader, data, funcMap, templateAwareFuncGenerators, nil)
}

// ProcessTemplateWithBaseTemplate processes the template using the specified data and a base template.
// The base template can contain common templates that can be referenced from the main template.
func ProcessTemplateWithBaseTemplate(reader io.Reader, data interface{}, funcMap template.FuncMap, templateAwareFuncGenerators TemplateAwareFuncMap, baseTemplate *template.Template) (io.Reader, error) {
	byteContent, err := readAll(reader)
	if err != nil {
		return nil, err
	}

	content, err := applyTemplateWithBase(string(byteContent), data, funcMap, templateAwareFuncGenerators, baseTemplate)
	if err != nil {
		log.Println("Error while generating output file from template")
		return nil, err
	}

	return content, nil
}
