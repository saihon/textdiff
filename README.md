## textdiff

for debug

<br/>

## example

<br/>

```go

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/saihon/textdiff"
)

func main() {
	text1 := `
foo
bar
baz
`
	text2 := `
foo
bar
vaz
`

	t := textdiff.New(strings.NewReader(text1), strings.NewReader(text2))

	for d := range t.Scan() {
		fmt.Fprintf(os.Stderr,
			"%d-%d:\n  TEXT-1: %s\n  TEXT-2: %s\n",
			d.Line, d.Index, d.Text1, d.Text2)
	}

	if err := t.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}


```
<br/>
<br/>
