package traceback

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type PrintFormat int

const (
	Text PrintFormat = iota
	Column
	Quickfix
	JSON
)

type PrintConfig struct {
	Format               PrintFormat
	PreserveSourcePrefix bool
}

func Fprint(out io.Writer, trace *Traceback, config PrintConfig) {
	switch config.Format {
	case Text:
		printText(out, trace, config)
	case Column:
		printColumn(out, trace, config)
	case Quickfix:
		printQuickfix(out, trace, config)
	case JSON:
		printJSON(out, trace, config)
	default:
		panic("unknown format: " + strconv.Itoa(int(config.Format)))
	}
}

func printText(out io.Writer, trace *Traceback, config PrintConfig) {
	tr := *trace
	if !config.PreserveSourcePrefix {
		tr.Stacks = TrimSourcePrefix(trace.Stacks)
	}
	if tr.Reason != "" {
		fmt.Fprintln(out, tr.Reason)
		fmt.Fprintln(out)
	}
	for i, s := range tr.Stacks {
		if i != 0 {
			fmt.Fprintln(out)
		}
		fmt.Fprintf(out, "goroutine %d [%s]:\n", s.ID, s.Status)
		for _, c := range s.Calls {
			fmt.Fprintf(out, "%s\n", c.Func)
			fmt.Fprintf(out, "\t%s:%d\n", c.Source, c.Line)
		}
	}
}

func printColumn(out io.Writer, trace *Traceback, config PrintConfig) {
	tr := *trace
	if !config.PreserveSourcePrefix {
		tr.Stacks = TrimSourcePrefix(trace.Stacks)
	}
	maxwidth := int(0)
	for _, s := range tr.Stacks {
		for _, c := range s.Calls {
			if maxwidth < len(c.Func) {
				maxwidth = len(c.Func)
			}
		}
	}
	if tr.Reason != "" {
		fmt.Fprintln(out, tr.Reason)
		fmt.Fprintln(out)
	}
	for i, s := range tr.Stacks {
		if i != 0 {
			fmt.Fprintln(out)
		}
		fmt.Fprintf(out, "goroutine %d [%s]\n", s.ID, s.Status)
		for _, c := range s.Calls {
			dw := maxwidth - len(c.Func)
			fmt.Fprintf(out, "  %s%s  %s:%d\n", c.Func, strings.Repeat(" ", dw), c.Source, c.Line)
		}
	}
}

func printQuickfix(out io.Writer, trace *Traceback, config PrintConfig) {
	for _, s := range trace.Stacks {
		for _, c := range s.Calls {
			msg := fmt.Sprintf("goroutine %d [%s]", s.ID, s.Status)
			fmt.Fprintf(out, "%s:%d: %s\n", c.Source, c.Line, msg)
		}
	}
}

func printJSON(out io.Writer, trace *Traceback, config PrintConfig) {
	tr := *trace
	tr.Stacks = TrimSourcePrefix(trace.Stacks)
	b, err := json.MarshalIndent(&tr, "", "\t")
	if err != nil {
		panic(err)
	}
	out.Write(b)
}
