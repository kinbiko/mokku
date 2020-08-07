package mokku

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTemplate(t *testing.T) {
	backupValue := os.Getenv(mokkuTemplatePathEnvName)

	for _, tc := range []struct {
		name     string
		in       *targetInterface
		exp      string
		tearUp   func(*testing.T)
		tearDown func(*testing.T)
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
	NoReturnParamFunc func( a string )
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
func (m *FooBarMock) NoReturnParam( a string ) {
	if m.NoReturnParamFunc == nil {
		panic("unexpected call to NoReturnParam")
	}
	m.NoReturnParamFunc( a )
}
`,

			in: &targetInterface{
				TypeName: "FooBar",
				Methods: []*method{
					{Name: "Act", Signature: "( ) error", OrderedParams: "( )", HasReturn: true},
					{Name: "DoStuff", Signature: "( a , b string, other ... interface{} ) ( int , error )", OrderedParams: "( a , b , other ... )", HasReturn: true},
					{Name: "NoReturnParam", Signature: "( a string )", OrderedParams: "( a )", HasReturn: false},
				},
			},
		},

		{
			name: "from template",
			exp: `
type BazMock struct { 
	HelloWorldFunc func( ) error
}
`,

			in: &targetInterface{
				TypeName: "Baz",
				Methods: []*method{
					{Name: "HelloWorld", Signature: "( ) error", OrderedParams: "( )", HasReturn: true},
				},
			},

			tearUp: func(t *testing.T) {
				err := os.Setenv(mokkuTemplatePathEnvName, filepath.Join("testdata", "test.tpl"))
				if err != nil {
					t.Fatal(err)
				}
			},

			tearDown: func(t *testing.T) {
				err := os.Setenv(mokkuTemplatePathEnvName, backupValue)
				if err != nil {
					t.Fatal(err)
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.tearUp != nil {
				tc.tearUp(t)
			}
			if tc.tearDown != nil {
				defer tc.tearDown(t)
			}

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
