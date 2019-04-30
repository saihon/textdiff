package textdiff

import (
	"bufio"
	"errors"
	"io"
	"unicode/utf8"
)

func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// find the next line feed index.
	n := -1
	for i, c := range data {
		if c == '\r' || c == '\n' {
			n = i
			break
		}
	}

	if n >= 0 {
		if data[n] == '\r' {
			// in the case of CR/LF, to delete the "\n" after "\r".
			if len(data) > n+1 && data[n+1] == '\n' {
				return n + 2, data[0:n], nil
			}
		}
		return n + 1, data[0:n], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

type TextDiff struct {
	scanner1        *bufio.Scanner
	scanner2        *bufio.Scanner
	StopImmediately bool
}

func New(r1, r2 io.Reader) *TextDiff {
	t := &TextDiff{
		scanner1: bufio.NewScanner(r1),
		scanner2: bufio.NewScanner(r2),
	}
	t.scanner1.Split(scanLines)
	t.scanner2.Split(scanLines)
	return t
}

func (t *TextDiff) Err() error {
	s := ""
	if err := t.scanner1.Err(); err != nil {
		s += err.Error()
	}
	if err := t.scanner2.Err(); err != nil {
		if len(s) > 0 {
			s += " "
		}
		s += err.Error()
	}
	if len(s) > 0 {
		return errors.New(s)
	}
	return nil
}

func (t *TextDiff) diff() int {
	s1 := t.scanner1.Text()
	s2 := t.scanner2.Text()
	index := 0
	for i, r := range s1 {
		if i > len(s2)-1 {
			return index
		}

		if r < utf8.RuneSelf {
			if r != rune(s2[i]) {
				return index
			}
		} else {
			c, _ := utf8.DecodeRuneInString(s2[i:])
			if r != c {
				return index
			}
		}
		index++
	}

	return -1
}

type Diff struct {
	Line  int
	Index int
	Text1 string
	Text2 string
}

func (t *TextDiff) scan(ch chan<- *Diff) {
	for i := 0; ; i++ {
		ok1 := t.scanner1.Scan()
		ok2 := t.scanner2.Scan()
		if !ok1 || !ok2 {
			switch {
			case ok1 && !ok2:
				ch <- &Diff{Line: i, Text1: t.scanner1.Text()}
			case !ok1 && ok2:
				ch <- &Diff{Line: i, Text2: t.scanner2.Text()}
			}
			break
		}

		if n := t.diff(); n != -1 {
			ch <- &Diff{
				Line:  i,
				Index: n,
				Text1: t.scanner1.Text(),
				Text2: t.scanner2.Text(),
			}

			if t.StopImmediately {
				break
			}
		}
	}
	close(ch)
}

// Scan
func (t *TextDiff) Scan() <-chan *Diff {
	ch := make(chan *Diff)
	go t.scan(ch)
	return ch
}
