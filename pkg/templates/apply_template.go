package templates

import (
	"bytes"
	"io"
	"text/template"
)

type TemplateAwareFuncMap map[string]func(*template.Template) any

func applyTemplate(templateContent string, config interface{}, funcMap template.FuncMap, templateAwareFuncGenerators TemplateAwareFuncMap) (io.Reader, error) {
	tpl := template.New("")

	templateAwareFuncMap := make(template.FuncMap, len(templateAwareFuncGenerators))
	for fnName, fnGen := range templateAwareFuncGenerators {
		templateAwareFuncMap[fnName] = fnGen(tpl)
	}

	tpl = tpl.Funcs(funcMap).Funcs(templateAwareFuncMap)

	tpl, err := tpl.Parse(templateContent)
	if err != nil {
		return nil, err
	}

	var result bytes.Buffer
	err = tpl.Execute(&result, config)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
