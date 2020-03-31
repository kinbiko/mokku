package mokku

import (
	"errors"
	"fmt"
	"go/scanner"
	"go/token"
	"strings"
)

type targetInterface struct {
	TypeName string
	Methods  []*method
}

type method struct {
	Name string

	// e.g. '( foo , bar string ) ( a int , err error )'
	// extracting it all as a string to keep things simple while still being
	// able flexible.
	Signature string

	// e.g. '( foo , bar )', used for passing parameters from the mock's method
	// to the mock struct's func property
	OrderedParams string

	HasReturn bool
}

func (m *method) String() string {
	return fmt.Sprintf("<Name: %s, Signature: %s>", m.Name, m.Signature)
}

type parser struct {
	s   *scanner.Scanner
	src []byte
}

func newParser(src []byte) *parser {
	var (
		fs = token.NewFileSet()
		s  = &scanner.Scanner{}
	)
	// currently the only scanner.Mode option is to include
	// comments (and a private option for not including semicolons).
	// Selecting 0 to *not* scan comments.
	bareEssentials := scanner.Mode(0)
	s.Init(fs.AddFile("", fs.Base(), len(src)), src, nil, bareEssentials)
	return &parser{s: s, src: src}
}

func (p *parser) parse() (*targetInterface, error) {
	name, err := p.lookForItfcName()
	if err != nil {
		return nil, fmt.Errorf("unable to find interface name in '%s':%w", p.src, err)
	}

	methods, err := p.lookForMethods()
	if err != nil {
		return nil, err
	}

	return &targetInterface{TypeName: name, Methods: methods}, nil
}

func (p *parser) lookForItfcName() (string, error) {
	for {
		tok, _ := p.scan(1)
		switch tok {
		case token.EOF:
			return "", errors.New("unable to find interface declaration")
		case token.INTERFACE:
			return "", nil // Treat as an anonymous interface
		case token.TYPE: // <type> Foo interface
			tok, lit := p.scan(2)
			if tok == token.IDENT { // type <Foo> interface
				_, tok, _ := p.s.Scan()
				if tok == token.INTERFACE { // type Foo <interface>
					return lit, nil
				}
			}
		}
	}
}

func (p *parser) lookForMethods() ([]*method, error) {
	for {
		tok, _ := p.scan(3)
		if tok == token.EOF {
			return nil, errors.New("unable to find method definitions for interface")
		}
		if tok == token.LBRACE {
			break
		}
	}

	methods := []*method{}

	for {
		tok, lit := p.scan(4)
		if tok == token.EOF {
			return nil, errors.New("unable to find method definition after '{'")
		}
		if tok == token.RBRACE {
			break
		}
		if tok == token.IDENT {
			methodName := lit
			m, err := p.lookForMethod(methodName)
			if err != nil {
				return nil, err
			}
			methods = append(methods, m)
		}
	}
	return methods, nil
}

func (p *parser) lookForMethod(methodName string) (*method, error) {
	collect := []string{}

	for {
		tok, lit := p.scan(5)
		switch tok {
		case token.COMMA,
			token.LPAREN,
			token.RPAREN,
			token.LBRACK,
			token.RBRACK,
			token.LBRACE,
			token.RBRACE,
			token.PERIOD,
			token.MUL, // pointer
			token.ELLIPSIS:
			collect = append(collect, tok.String())
		case token.EOF:
			return nil, fmt.Errorf("unable to parse method definition for '%s'", methodName)
		}

		if tok == token.SEMICOLON {
			break
		}

		if lit != "" {
			collect = append(collect, lit)
		}
	}

	// Joining with a " " in so that the tokeniser can differentiate between
	// "a string" and "astring". Granted, this isn't going to be nice to read,
	// but the idea is to run gofmt (or goimports) later down the line, to
	// enforce standard go syntax.
	sig := strings.Join(collect, " ")
	orderedParams := parseArgs(sig)
	// TODO: clean up this ugly hack
	ind := strings.Index(sig, ")")
	hasReturn := false
	if ind > 0 {
		hasReturn = len(sig[ind:]) > 2
	}

	return &method{
		Name:          methodName,
		Signature:     sig,
		OrderedParams: orderedParams,
		HasReturn:     hasReturn,
	}, nil
}

func parseArgs(src string) string {
	if src == "" {
		// in this case its likely a case of having an embedded interface, and
		// thus args don't make sense.
		return ""
	}

	sig := src[:strings.Index(src, ")")+1]
	if sig == "( )" {
		return sig
	}

	var (
		collect    = []string{}
		addNext    = true
		hasVarargs = false

		pp = newParser([]byte(sig))
	)

	for {
		tok, lit := pp.scan(6)
		if addNext {
			if lit == "" {
				lit = tok.String()
			}
			collect = append(collect, lit)
		}
		switch tok {
		case token.COMMA:
			collect = append(collect, tok.String())
			addNext = true
		case token.ELLIPSIS:
			hasVarargs = true
		case token.LPAREN:
			addNext = true
		case token.EOF:
			if hasVarargs {
				collect = append(collect, "...")
			}
			return strings.Join(append(collect, ")"), " ")
		default:
			addNext = false
		}
	}
}

// scan reads the next token to scan.
// The input argument only serves as a code-base location for debugging, if
// necessary.
func (p *parser) scan(_ int) (token.Token, string) {
	_, tok, lit := p.s.Scan()
	// NOTE: this is a good spot to do fmt.Println debugging if needed
	return tok, lit
}
