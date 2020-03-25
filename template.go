package mokku

import (
	"bytes"
	"text/template"
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
	return m.{{$val.Name}}Func{{$val.OrderedParams}}
}
{{ end }}{{ end }}`

func mockFromTemplate(defn *targetInterface) ([]byte, error) {
	tmpl, err := template.New("mock").Parse(tpl)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, defn); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
