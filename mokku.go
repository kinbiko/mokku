package mokku

import "go/format"

// Config is currently ignored but is expected to contain various configuration
// options for this tool given by user-provided flags in the future.
type Config struct {
	// Intentionally empty at the moment.
	// Included only to avoid breaking backwards compatibility if a newer
	// version of the package supports new features
}

// Mock creates the sourcecode of a mock implementation of the interface
// sourcecode defined in the given byte array.
func Mock(_ Config, src []byte) ([]byte, error) {
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
