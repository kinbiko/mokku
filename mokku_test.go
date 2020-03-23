package mokku_test

import (
	"context"
	"testing"

	"github.com/kinbiko/mokku"
)

func TestIntegration(t *testing.T) {
	got, err := mokku.Mock(context.Background(), []byte(`type Foo interface{Act()}`))
	if err != nil {
		t.Fatalf("unexpected error '%s'", err.Error())
	}
	exp := `
// FooMock is a mock implementation of Foo.
type FooMock struct {
	ActImpl func()
}

func (f *FooMock) Act() {
	if f.ActImpl == nil {
		panic("unexpected call to Act")
	}
	f.ActImpl()
}`
	if string(got) != exp {
		t.Errorf("unexpected mock created:\n%s\n\nexpected:\n%s", got, exp)
	}
}
