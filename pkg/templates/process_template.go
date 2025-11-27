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

// ProcessTemplate processes the template using the specified data
func ProcessTemplateWithTemplateAware(reader io.Reader, data interface{}, funcMap template.FuncMap, templateAwareFuncGenerators TemplateAwareFuncMap) (io.Reader, error) {
	byteContent, err := readAll(reader)
	if err != nil {
		return nil, err
	}

	content, err := applyTemplate(string(byteContent), data, funcMap, templateAwareFuncGenerators)
	if err != nil {
		log.Println("Error while generating output file from template")
		return nil, err
	}

	return content, nil
}
