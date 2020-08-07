package mokku

import "go/format"

// Config is defines all configuration options for mokku.
// In particular, this package treats additions to this struct as *non-breaking* changes.
type Config struct {
	// TemplateStr is mock template string.
	// The template will attempt to fill in the following fields:
	//  TypeName -- string. E.g. "MyInterface"
	//	Methods  -- composite type containing:
	//		Name			-- string. E.g. "DoStuff"
	//		Signature		-- string. E.g. "(a, b int) error"
	//		OrderedParams	-- string. E.g. "(a, b)"
	//		HasReturn		-- bool. Identifies whether or not the mocked method should return anything.
	TemplateStr string
}

// Mock creates the sourcecode of a mock implementation of the interface
// sourcecode defined in the given byte array.
func Mock(config Config, src []byte) ([]byte, error) {
	target, err := newParser(src).parse()
	if err != nil {
		return nil, err
	}
	mft, err := mockFromTemplate(target, config.TemplateStr)
	if err != nil {
		return nil, err
	}
	return format.Source(mft)
}
