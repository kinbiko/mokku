package mokku

import (
	"testing"
)

func TestParser(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  string
		exp  *targetInterface
	}{
		{
			name: "named empty interface",
			src:  `type Foo interface{}`,
			exp:  &targetInterface{name: "Foo"},
		},
		{
			name: "anonymous empty interface",
			src:  `interface{}`,
			exp:  &targetInterface{},
		},
		{
			name: "single niladic method with no return params",
			src: `type Bar interface{
				Act()
			}`,
			exp: &targetInterface{
				name:    "Bar",
				methods: []*method{{"Act", "( )"}},
			},
		},
		{
			name: "two niladic methods with no return params",
			src: `type FooBar interface{
					Act()
					Do()
				}`,
			exp: &targetInterface{
				name:    "FooBar",
				methods: []*method{{"Act", "( )"}, {"Do", "( )"}},
			},
		},
		{
			name: "single method with single input parameter and no return params",
			src: `type FooBar interface{
					Act(x string)
				}`,
			exp: &targetInterface{
				name:    "FooBar",
				methods: []*method{{name: "Act", signature: `( x string )`}},
			},
		},
		{
			name: "single method with multiple complex input parameters and no return params",
			src: `type FooBar interface{
					Act(x, y string, z chan []struct{a [0]int})
				}`,
			exp: &targetInterface{
				name:    "FooBar",
				methods: []*method{{name: "Act", signature: `( x , y string , z chan [ ] struct { a [ 0 ] int } )`}},
			},
		},

		{
			name: "single method with no input parameters and one return parameter",
			src: `type FooBar interface {
				Act() error
			}`,
			exp: &targetInterface{
				name:    "FooBar",
				methods: []*method{{name: "Act", signature: `( ) error`}},
			},
		},

		{
			name: "mega complex example",
			src: `type GoodLuck interface {
				First()
				Second(ctx context.Context, f []*fish, s, ss string) (error, int, chan struct{})
				Third() (a, b string, e error)
			}`,
			exp: &targetInterface{
				name: "GoodLuck",
				methods: []*method{
					{name: "First", signature: `( )`},
					{name: "Second", signature: `( ctx context . Context , f [ ] fish , s , ss string ) ( error , int , chan struct { } )`},
					{name: "Third", signature: `( ) ( a , b string , e error )`},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := newParser([]byte(tc.src))

			got, err := p.parse()
			if err != nil {
				t.Fatalf("got unexpected error when calling parse parser: '%s'", err.Error())
			} else if got == nil {
				t.Fatal("expected non-nil parsed result but was nil")
			}

			if got.name != tc.exp.name {
				t.Errorf("expected name '%s' but got '%s'", tc.exp.name, got.name)
			}

			for i := range tc.exp.methods {
				expLen, gotLen := len(tc.exp.methods), len(got.methods)
				if expLen != gotLen {
					t.Errorf("expected %d methods but got %d", expLen, gotLen)
					t.Errorf("exp: '%+v'", tc.exp.methods)
					t.Errorf("got: '%+v'", got.methods)
					break
				}

				expName, gotName := tc.exp.methods[i].name, got.methods[i].name
				if expName != gotName {
					t.Errorf("expected method of index %d to have name '%s' but was '%s'", i, expName, gotName)
				}
				expSignature, gotSignature := tc.exp.methods[i].signature, got.methods[i].signature
				if expSignature != gotSignature {
					t.Errorf("expected method of index %d to had different signatures:\nexp: %s\nwas: %s", i, expSignature, gotSignature)
				}
			}
		})
	}
}
