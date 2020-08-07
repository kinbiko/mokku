package mokku

import (
	"bytes"
	"text/template"
)

func mockFromTemplate(defn *targetInterface, templateStr string) ([]byte, error) {
	tmpl, err := template.New("mock").Parse(templateStr)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, defn); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
