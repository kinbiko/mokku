package templates

func init() {
	templateMap["default"] = defaultTemplate
}

const defaultTemplate = `
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
