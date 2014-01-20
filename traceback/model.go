// This is a derivative work of roger peppe's stackgraph command.
// For more details, see http://code.google.com/p/rog-go/
//

package traceback

import (
	"fmt"
	"io"
	"strings"
)

type StackStatus string

const (
	StackStatusChanReceive = "chan receive"
	StackStatusSemAcquire  = "semacquire"
	StackStatusRunning     = "running"
	StackStatusRunnable    = "runnable"
)

// Call represents a function call.
type Call struct {
	Func   string
	Source string
	Line   int
	Args   []uint64
}

// Stack represents the call stack of a goroutine.
type Stack struct {
	ID     int
	Status StackStatus
	Calls  []Call
}

func Fprint(out io.Writer, stacks []*Stack) {
	for i, s := range stacks {
		if i != 0 {
			fmt.Fprintln(out)
		}
		fmt.Fprintf(out, "goroutine %d [%s]\n", s.ID, s.Status)
		maxwidth := int(0)
		for _, c := range s.Calls {
			if maxwidth < len(c.Func) {
				maxwidth = len(c.Func)
			}
		}
		for _, c := range s.Calls {
			dw := maxwidth - len(c.Func)
			fmt.Fprintf(out, "  %s%s %s:%d\n", c.Func, strings.Repeat(" ", dw), c.Source, c.Line)
		}
	}
}
