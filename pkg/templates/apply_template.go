package templates

import (
	"bytes"
	"io"
	"text/template"
)

func applyTemplate(templateContent string, config interface{}, funcMap template.FuncMap) (io.Reader, error) {
	template, err := template.New("").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return nil, err
	}

	var result bytes.Buffer
	err = template.Execute(&result, config)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
