package mokku

import (
	"testing"

	"github.com/kinbiko/mokku/templates"
)

func TestTemplate(t *testing.T) {
	defaultTemplate := templates.GetDefault()
	type args struct {
		defn *targetInterface
		templateStr string
	}
	for _, tc := range []struct {
		name string
		args args
		exp  string
	}{
		{
			name: "basic case",
			exp: `
type Mock struct { 
}
`,
			args: args{
				defn: &targetInterface{},
				templateStr: defaultTemplate,
			},
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

			args: args{
				defn: &targetInterface{
					TypeName: "FooBar",
					Methods: []*method{
						{Name: "Act", Signature: "( ) error", OrderedParams: "( )", HasReturn: true},
						{Name: "DoStuff", Signature: "( a , b string, other ... interface{} ) ( int , error )", OrderedParams: "( a , b , other ... )", HasReturn: true},
						{Name: "NoReturnParam", Signature: "( a string )", OrderedParams: "( a )", HasReturn: false},
					},
				},
				templateStr: defaultTemplate,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, err := mockFromTemplate(tc.args.defn, tc.args.templateStr)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}

			if got := string(b); got != tc.exp {
				t.Errorf("got different output than what was expected:\n%s", got)
			}
		})
	}
}
