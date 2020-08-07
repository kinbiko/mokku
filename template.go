package mokku

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"
)

const (
	mokkuTemplatePathEnvName = "MOKKU_TEMPLATE_PATH"
)

const tpl = `
type {{.TypeName}}Mock struct { {{ range .Methods }}
	{{.Name}}Func func{{.Signature}}{{ end }}
}
{{if .Methods }}{{$typeName := .TypeName}}
{{range $val := .Methods}}func (m *{{$typeName}}Mock) {{$val.Name}}{{$val.Signature}} {
	if m.{{$val.Name}}Func == nil {
		panic("unexpected call to {{$val.Name}}")
	}
	{{if $val.HasReturn}}return {{ end }}m.{{$val.Name}}Func{{$val.OrderedParams}}
}
{{ end }}{{ end }}`

func mockFromTemplate(defn *targetInterface) ([]byte, error) {
	// set default template to source template
	src := tpl

	p := os.Getenv(mokkuTemplatePathEnvName)
	if p != "" {
		// if succeeding to read from the file,
		// set to source template
		if b, err := ioutil.ReadFile(p); err == nil {
			src = string(b)
		}
	}

	tmpl, err := template.New("mock").Parse(src)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, defn); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
