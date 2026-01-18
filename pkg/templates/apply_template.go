package templates

import (
	"bytes"
	"io"
	"text/template"
)

type TemplateAwareFuncMap map[string]func(*template.Template) any

func applyTemplate(templateContent string, config interface{}, funcMap template.FuncMap, templateAwareFuncGenerators TemplateAwareFuncMap) (io.Reader, error) {
	return applyTemplateWithBase(templateContent, config, funcMap, templateAwareFuncGenerators, nil)
}

func applyTemplateWithBase(templateContent string, config interface{}, funcMap template.FuncMap, templateAwareFuncGenerators TemplateAwareFuncMap, baseTemplate *template.Template) (io.Reader, error) {
	var tpl *template.Template
	if baseTemplate != nil {
		// Clone the base template to inherit all associated templates (common templates)
		var err error
		tpl, err = baseTemplate.Clone()
		if err != nil {
			return nil, err
		}
		// Create a new template within the cloned base to parse the new content
		tpl = tpl.New("")
	} else {
		tpl = template.New("")
	}

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
