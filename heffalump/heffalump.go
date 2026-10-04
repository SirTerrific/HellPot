/*
Package heffalump attempts to encapsulate the original work by carlmjohnson on heffalump
https://github.com/carlmjohnson/heffalump
*/
package heffalump

import (
	"bufio"
	"io"

	"github.com/SirTerrific/HellPot/internal/config"
)

var log = config.GetLogger()

// DefaultHeffalump represents a Heffalump type
var DefaultHeffalump *Heffalump

// Heffalump represents our markov map from Heffalump
type Heffalump struct {
	mm MarkovMap
}

// NewHeffalump instantiates a new Heffalump for markov generation and io operations.
// buffsize is kept for API compatibility but unused: WriteHell writes to a *bufio.Writer,
// which implements io.ReaderFrom, so io.CopyBuffer never touched the pooled buffer.
func NewHeffalump(mm MarkovMap, buffsize int) *Heffalump {
	return &Heffalump{mm: mm}
}

// WriteHell writes markov chain heffalump hell to the provided io.Writer
func (h *Heffalump) WriteHell(bw *bufio.Writer) (int64, error) {
	var n int64
	var err error

	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("caller", r).Msg("panic recovered!")
		}
	}()

	if _, err = bw.WriteString("<html>\n<body>\n"); err != nil {
		return n, err
	}
	n, _ = io.Copy(bw, h.mm)
	return n, nil
}
