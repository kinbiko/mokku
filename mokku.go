package mokku

import "go/format"

// Mock creates the sourcecode of a mock implementation of the interface
// sourcecode defined in the given byte array.
func Mock(src []byte) ([]byte, error) {
	target, err := newParser(src).parse()
	if err != nil {
		return nil, err
	}
	mft, err := mockFromTemplate(target)
	if err != nil {
		return nil, err
	}
	return format.Source(mft)
}
