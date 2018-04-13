// parse package contains registry of known parsers and allows to invoke them
// by name. Parser implementation should register a function of type parseFunc.
package parse

import (
	"fmt"
	"io"
	"log"

	"pass2pass/model"
)

// ParseAsync asynchronously parse input using parser of requested parser type.
// parserType must be a registered parser type.
func ParseAsync(
	parserType string,
	r io.Reader) (<-chan model.Entry, <-chan error) {
	p, ok := parsers[parserType]
	if !ok {
		errCh := make(chan error, 1)
		defer close(errCh)
		errCh <- model.NewFatalError(fmt.Errorf("unknown parser: %s", parserType))
		return nil, errCh
	}
	log.Printf("found parser function for parser type: %s", parserType)
	return p(r)
}

type parseFunc func(io.Reader) (<-chan model.Entry, <-chan error)

var parsers = map[string]parseFunc{}
