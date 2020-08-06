package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

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
	flag.Usage = func() { fmt.Println(usage) }
	templateName := flag.String("t", "", "TemplateFile")
	flag.Parse()

	templateStr := defaultTemplate
	if *templateName != "" {
		var err error
		templateStr, err = loadTemplate(*templateName)
		if err != nil {
			errorOut(err)
		}
	}

	s, err := clipboard.ReadAll()
	if err != nil {
		errorOut(err)
	}

	mock, err := mokku.Mock(mokku.Config{TemplateStr: templateStr}, []byte(s))
	if err != nil {
		errorOut(err)
	}

	if err = clipboard.WriteAll(string(mock)); err != nil {
		errorOut(err)
	}
}

// loadTemplate loads template string from filePath
func loadTemplate(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to read file %s: %v", filePath, err))
	}
	templateStr := string(content)
	return templateStr, nil
}

func errorOut(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	flag.Usage()
	os.Exit(1)
}
