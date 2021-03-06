package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/t-yuki/gotfmt/traceback"
)

func Convert(in io.Reader, out io.Writer) *traceback.Traceback {
	var wr io.Writer = ioutil.Discard
	switch *format {
	case "raw":
		io.Copy(out, in)
		return nil
	case "text", "pretty":
		wr = out
	default:
		wr = ioutil.Discard
	}

	trace, err := traceback.ParseTraceback(in, wr)
	if err != nil {
		panic(err)
	}
	if trace == nil {
		return nil
	}

	filtered := ApplyFilters(trace)

	switch *format {
	case "text":
		if len(filtered.Stacks) != 0 {
			fmt.Fprintln(out)
		}
		traceback.Fprint(out, filtered, traceback.PrintConfig{Format: traceback.Text})
	case "pretty":
		if len(filtered.Stacks) != 0 {
			fmt.Fprintln(out)
		}
		traceback.Fprint(out, filtered, traceback.PrintConfig{Format: traceback.Pretty})
	case "qfix":
		traceback.Fprint(out, filtered, traceback.PrintConfig{Format: traceback.Quickfix})
	case "json":
		traceback.Fprint(out, filtered, traceback.PrintConfig{Format: traceback.JSON})
	default:
		panic("unsupported format: " + *format)
	}
	return trace
}

func ApplyFilters(src *traceback.Traceback) *traceback.Traceback {
	trace := src.Clone()
	stacks := trace.Stacks
	if strings.Contains(*filter, "notest") {
		stacks = traceback.ExcludeGotest(stacks)
	}
	if strings.Contains(*filter, "nostd") {
		stacks = traceback.ExcludeGoroot(stacks, false)
	} else if strings.Contains(*filter, "trimstd") {
		stacks = traceback.ExcludeGoroot(stacks, true)
	}
	if strings.Contains(*filter, "top") {
		stacks = traceback.ExcludeLowers(stacks)
	}
	trace.Stacks = stacks
	return &trace
}
