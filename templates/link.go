package templates

import (
	"errors"
	"fmt"
)

var (
	// You can add your template to this Map.
	// See toy_template.go
	templateMap = make(map[string]string)
)

// Get returns template string from template name
func Get(templateName string) (string, error) {
	if tmplStr, ok := templateMap[templateName]; ok {
		return tmplStr, nil
	}
	return "", errors.New(fmt.Sprintf("unable to find template: %s", templateName))
}

// Get returns default template string
func GetDefault() (string) {
	return defaultTemplate
}
