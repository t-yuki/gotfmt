package traceback

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/t-yuki/gotfmt/thirdparty/go-runewidth"
)

type PrintFormat int

const (
	Text PrintFormat = iota
	Pretty
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
	case Pretty:
		printPretty(out, trace, config)
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

var gowidth = runewidth.Condition{false}

func formatSourceLine(c Call) (line string, width int) {
	line = fmt.Sprintf("%s:%d", c.Source, c.Line)
	width = gowidth.StringWidth(line)
	return
}

func printPretty(out io.Writer, trace *Traceback, config PrintConfig) {
	tr := *trace
	if !config.PreserveSourcePrefix {
		tr.Stacks = TrimSourcePrefix(trace.Stacks)
	}
	if tr.Reason != "" {
		fmt.Fprintln(out, tr.Reason)
		fmt.Fprintln(out)
	}
	for i, s := range tr.Stacks {
		call2width := make([]int, len(s.Calls))
		width2call := make(map[int][]int)
		minwidth, maxwidth := 0, 0
		for i, c := range s.Calls {
			_, width := formatSourceLine(c)
			width2call[width] = append(width2call[width], i)
			if minwidth == 0 || minwidth > width {
				minwidth = width
			}
			if maxwidth < width {
				maxwidth = width
			}
		}
		for w := minwidth; w <= maxwidth; {
			if len(width2call[w]) == 0 {
				w++
				continue
			}
			max2 := 0
			for v := w; v <= maxwidth && v < w+8; v++ {
				if len(width2call[v]) != 0 {
					max2 = v
				}
			}
			for v := w; v <= max2; v++ {
				calls := width2call[v]
				for _, c := range calls {
					dw := max2 - v
					call2width[c] = dw
				}
			}
			w = max2 + 1
		}

		if i != 0 {
			fmt.Fprintln(out)
		}
		fmt.Fprintf(out, "goroutine %d [%s]\n", s.ID, s.Status)
		for i, c := range s.Calls {
			src, _ := formatSourceLine(c)
			dw := call2width[i]

			fn := c.Func
			if idx := strings.LastIndex(fn, "/"); idx != -1 && idx+1 < len(fn) {
				fn = fn[idx+1:]
			}
			if idx := strings.Index(fn, "."); idx != -1 && idx+1 < len(fn) {
				fn = fn[idx+1:]
			}
			fmt.Fprintf(out, "  %s%s %s()\n", src, strings.Repeat(" ", dw), fn)
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
