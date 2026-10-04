package heffalump

import (
	"bufio"
	"errors"
	"strings"
	"testing"
)

// failAfter accepts limit bytes then fails, like a client hanging up.
type failAfter struct {
	sb    strings.Builder
	limit int
}

func (f *failAfter) Write(p []byte) (int, error) {
	if f.sb.Len()+len(p) > f.limit {
		return 0, errors.New("connection closed")
	}
	return f.sb.Write(p)
}

// Mirrors the router loop: WriteHell must stream markov output and the loop must end once the client is gone.
func TestWriteHellEndsOnClientError(t *testing.T) {
	w := &failAfter{limit: 64 << 10}
	bw := bufio.NewWriter(w)

	var total int64
	for i := 0; ; i++ {
		if i > 1000 {
			t.Fatal("WriteHell never reported the writer error")
		}
		n, err := DefaultHeffalump.WriteHell(bw)
		total += n
		if err != nil {
			break
		}
	}

	if total == 0 || !strings.HasPrefix(w.sb.String(), "<html>\n<body>\n") {
		t.Fatalf("unexpected output: total=%d prefix=%q", total, w.sb.String()[:min(20, w.sb.Len())])
	}
}
