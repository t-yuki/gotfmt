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
	StackStatusChanReceive   = "chan receive"
	StackStatusSemAcquire    = "semacquire"
	StackStatusRunning       = "running"
	StackStatusRunnable      = "runnable"
	StackStatusSleep         = "sleep"
	StackStatusFinalizerWait = "finalizer wait"
	StackStatusSyscall       = "syscall"
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

type PrintConfig struct {
	Quickfix      bool
	OmitGoroutine bool
}

func Fprint(out io.Writer, stacks []*Stack, config PrintConfig) {
	maxwidth := int(0)
	for _, s := range stacks {
		for _, c := range s.Calls {
			if maxwidth < len(c.Func) {
				maxwidth = len(c.Func)
			}
		}
	}
	for i, s := range stacks {
		if i != 0 {
			fmt.Fprintln(out)
		}
		if !config.OmitGoroutine && !config.Quickfix {
			fmt.Fprintf(out, "goroutine %d [%s]\n", s.ID, s.Status)
		}
		for _, c := range s.Calls {
			dw := maxwidth - len(c.Func)
			if config.Quickfix {
				msg := fmt.Sprintf("goroutine %d [%s]\n", s.ID, s.Status)
				fmt.Fprintf(out, "%s:%d: %s\n", c.Source, c.Line, msg)
			} else {
				fmt.Fprintf(out, "  %s%s  %s:%d\n", c.Func, strings.Repeat(" ", dw), c.Source, c.Line)
			}
		}
	}
}
