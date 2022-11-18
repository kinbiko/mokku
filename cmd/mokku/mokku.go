package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/kinbiko/mokku"
)

const usage = `Usage:
1. Copy the interface you want to mock
2. Run 'mokku'
3. Paste the mocked implementation that has been written to your clipboard`

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

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		flag.Usage()
		os.Exit(1)
	}
}

func run() error {
	flag.Usage = func() { fmt.Println(usage) }
	flag.Parse()

	templateStr, err := loadTemplate(os.Getenv("MOKKU_TEMPLATE_PATH"))
	if err != nil {
		return err
	}
	if templateStr == "" {
		templateStr = defaultTemplate
	}

	s, err := clipboard.ReadAll()
	if err != nil {
		return err
	}

	mock, err := mokku.Mock(mokku.Config{TemplateStr: templateStr}, []byte(s))
	if err != nil {
		return err
	}

	if err := clipboard.WriteAll(string(mock)); err != nil {
		return err
	}
	return nil
}

// loadTemplate loads template string from filePath, if there is one.
func loadTemplate(filePath string) (string, error) {
	if filePath == "" { // There may not be an external template path given
		return "", nil
	}
	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %v", filePath, err)
	}
	return string(content), nil
}
