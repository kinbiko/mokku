package mokku

import (
	"strings"
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
				methods: []*method{{"Act", "( )", "( )"}},
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
				methods: []*method{{"Act", "( )", "( )"}, {"Do", "( )", "( )"}},
			},
		},
		{
			name: "single method with single input parameter and no return params",
			src: `type FooBar interface{
					Act(x string)
				}`,
			exp: &targetInterface{
				name:    "FooBar",
				methods: []*method{{"Act", `( x string )`, "( x )"}},
			},
		},
		{
			name: "single method with multiple complex input parameters and no return params",
			src: `type FooBar interface{
					Act(x, y string, z chan []struct{a [0]int})
				}`,
			exp: &targetInterface{
				name:    "FooBar",
				methods: []*method{{"Act", `( x , y string , z chan [ ] struct { a [ 0 ] int } )`, "( x , y , z )"}},
			},
		},

		{
			name: "single method with no input parameters and one return parameter",
			src: `type FooBar interface {
				Act() error
			}`,
			exp: &targetInterface{
				name:    "FooBar",
				methods: []*method{{"Act", `( ) error`, "( )"}},
			},
		},

		{
			name: "mega complex example",
			src: `type GoodLuck interface {
				First()
				Second(ctx context.Context, _ []*fish, s, ss string) (error, int, chan struct{})
				Third(vararg ...map[string]interface{}) (a, b string, e error)
			}`,
			exp: &targetInterface{
				name: "GoodLuck",
				methods: []*method{
					{"First", `( )`, "( )"},
					{"Second", `( ctx context . Context , _ [ ] fish , s , ss string ) ( error , int , chan struct { } )`, "( ctx , _ , s , ss )"}, // TODO: figure out what the default value is likely to be for _s as _ is useless as params.
					{"Third", `( vararg ... map [ string ] interface { } ) ( a , b string , e error )`, "( vararg ... )"},
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
				expParams, gotParams := tc.exp.methods[i].orderedParams, got.methods[i].orderedParams
				if expParams != gotParams {
					t.Errorf("expected method of index %d to had different ordered params:\nexp: %s\nwas: %s", i, expParams, gotParams)
				}
			}
		})
	}

	t.Run("error cases", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			src  string
			exp  string
		}{
			{
				name: "just curly brackets",
				src:  `{}`,
				exp:  "unable to find interface declaration",
			},
			{
				name: "missing closing curly bracket",
				src:  `type Foo interface{`,
				exp:  "unable to find method definition",
			},
			{
				name: "missing closing round bracket",
				src:  `type Foo interface{ Act(a string, }`,
				exp:  "unable to find method definition",
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				p := newParser([]byte(tc.src))

				_, err := p.parse()
				if err == nil {
					t.Fatalf("missing expected error when calling parse parser with: '%s'", tc.src)
				}
				if !strings.Contains(err.Error(), tc.exp) {
					t.Errorf("expected an error message containing '%s' but got '%s'", tc.exp, err.Error())
				}
			})
		}
	})
}

func TestOrderedParamsFromSignature(t *testing.T) {
	for _, tc := range []struct{ name, sig, exp string }{
		{"not a method", "", ""},
		{"niladic", "( )", "( )"},
		{"niladic with return param", "( ) error", "( )"},
		{"niladic with return params", "( ) ( int, error )", "( )"},
		{"monadic", "( a string )", "( a )"},
		{"dyadic with first parameter type absent", "( a, b string )", "( a , b )"},
		{"dyadic", "( a int, b string )", "( a , b )"},
		{"triadic with only last parameter type present", "( a, b, c string )", "( a , b , c )"},
		{"complex case", "( a map[int]interface{}, b chan []struct{}, c ...string )", "( a , b , c ... )"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseArgs(tc.sig); got != tc.exp {
				t.Errorf("expected '%s' but got '%s'", tc.exp, got)
			}
		})
	}
}
