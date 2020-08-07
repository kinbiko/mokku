package mokku_test

import (
	"testing"

	"github.com/kinbiko/mokku"
)

func TestIntegration(t *testing.T) {
	const templateStr = `
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

	got, err := mokku.Mock(mokku.Config{TemplateStr: templateStr}, []byte(`
    type Foo interface {
		Act()
	}`))
	if err != nil {
		t.Fatalf("unexpected error '%s'", err.Error())
	}

	exp := `
type FooMock struct {
	ActFunc func()
}

func (m *FooMock) Act() {
	if m.ActFunc == nil {
		panic("unexpected call to Act")
	}
	m.ActFunc()
}
`
	if string(got) != exp {
		t.Errorf("unexpected mock created:\n%s\n\nexpected:\n%s", got, exp)
	}
}
