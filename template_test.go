package mokku

import "testing"

func TestTemplate(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   *targetInterface
		exp  string
	}{
		{
			name: "basic case",
			exp: `
type Mock struct { 
}
`,
			in: &targetInterface{},
		},

		{
			name: "advanced case",
			exp: `
type FooBarMock struct { 
	ActFunc func( ) error
	DoStuffFunc func( a , b string, other ... interface{} ) ( int , error )
}

func (m *FooBarMock) Act( ) error {
	if m.ActFunc == nil {
		panic("unexpected call to Act")
	}
	return m.ActFunc( )
}
func (m *FooBarMock) DoStuff( a , b string, other ... interface{} ) ( int , error ) {
	if m.DoStuffFunc == nil {
		panic("unexpected call to DoStuff")
	}
	return m.DoStuffFunc( a , b , other ... )
}
`,

			in: &targetInterface{
				TypeName: "FooBar",
				Methods: []*method{
					{Name: "Act", Signature: "( ) error", OrderedParams: "( )"},
					{Name: "DoStuff", Signature: "( a , b string, other ... interface{} ) ( int , error )", OrderedParams: "( a , b , other ... )"},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, err := mockFromTemplate(tc.in)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}

			if got := string(b); got != tc.exp {
				t.Errorf("got different output than what was expected:\n%s", got)
			}
		})
	}
}
