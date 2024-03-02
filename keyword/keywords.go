package keyword

import (
	"bufio"
	"io"
	"strings"
)

type Keywords []Keyword

func (ks Keywords) Match(r io.Reader) (Keyword, bool) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		for _, k := range ks {
			if strings.Contains(line, k.Keyword) {
				return k, true
			}
		}
	}
	return Keyword{}, false
}

type Keyword struct {
	Keyword string
	Color   string
}
