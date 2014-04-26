package traceback

import (
	"fmt"
	"io"
	"strings"
)

type PrintConfig struct {
	Quickfix      bool
	OmitGoroutine bool
}

func Fprint(out io.Writer, trace *Traceback, config PrintConfig) {
	maxwidth := int(0)
	for _, s := range trace.Stacks {
		for _, c := range s.Calls {
			if maxwidth < len(c.Func) {
				maxwidth = len(c.Func)
			}
		}
	}
	if !config.Quickfix && trace.Reason != "" {
		fmt.Fprintln(out, trace.Reason)
		fmt.Fprintln(out)
	}
	for i, s := range trace.Stacks {
		if i != 0 && !config.Quickfix {
			fmt.Fprintln(out)
		}
		if !config.OmitGoroutine && !config.Quickfix {
			fmt.Fprintf(out, "goroutine %d [%s]\n", s.ID, s.Status)
		}
		for _, c := range s.Calls {
			dw := maxwidth - len(c.Func)
			if config.Quickfix {
				msg := fmt.Sprintf("goroutine %d [%s]", s.ID, s.Status)
				fmt.Fprintf(out, "%s:%d: %s\n", c.Source, c.Line, msg)
			} else {
				fmt.Fprintf(out, "  %s%s  %s:%d\n", c.Func, strings.Repeat(" ", dw), c.Source, c.Line)
			}
		}
	}
}
