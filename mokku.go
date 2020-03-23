package mokku

import "context"

func Mock(ctx context.Context, data []byte) ([]byte, error) {
	// TODO: hardcoded for now
	return []byte(`
// FooMock is a mock implementation of Foo.
type FooMock struct {
	ActImpl func()
}

func (f *FooMock) Act() {
	if f.ActImpl == nil {
		panic("unexpected call to Act")
	}
	f.ActImpl()
}`), nil
}
