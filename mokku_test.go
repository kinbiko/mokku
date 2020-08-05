package mokku_test

import (
	"testing"

	"github.com/kinbiko/mokku"
	"github.com/kinbiko/mokku/templates"
)

func TestIntegration(t *testing.T) {
	got, err := mokku.Mock(mokku.Config{TemplateStr: templates.GetDefault()}, []byte(`
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
