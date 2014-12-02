package cmd

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	"github.com/t-yuki/gotfmt/traceback"
)

func Convert(in io.Reader, out io.Writer) {
	var wr io.Writer = ioutil.Discard
	switch *format {
	case "raw":
		fallthrough
	default:
		io.Copy(out, in)
		return
	case "col":
		wr = out
	}
	trace, err := traceback.ParseTraceback(in, wr)
	if err != nil {
		panic(err)
	}
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
	switch *format {
	case "col":
		stacks = traceback.TrimSourcePrefix(stacks)
		trace.Stacks = stacks
		traceback.Fprint(out, trace, traceback.PrintConfig{})
	case "qfix":
		trace.Stacks = stacks
		traceback.Fprint(out, trace, traceback.PrintConfig{Quickfix: true})
	case "json":
		stacks = traceback.TrimSourcePrefix(stacks)
		trace.Stacks = stacks
		b, err := json.MarshalIndent(trace, "", "\t")
		if err != nil {
			panic(err)
		}
		out.Write(b)
	default:
		panic("unsupported format: " + *format)
	}
}
