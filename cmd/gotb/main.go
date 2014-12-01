package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/t-yuki/gotfmt/traceback"
)

var (
	filterGOROOT  = flag.String("filter.goroot", "chop", "set GOROOT filter mode. none: show all GOROOT functions. all: exclude all GOROOT functions completely. chop: omit consecutive GOROOT functions from the top of stack except one")
	filterTopOnly = flag.Bool("filter.toponly", false, "print each top function of stacks")
	format        = flag.String("out.format", "column", "column: two column. qfix: vim quickfix list format. it should use with `toponly`")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "gotb - Go Trace Back formatter")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, syscall.SIGQUIT)
	convert(os.Stdin, os.Stdout)
}

func convert(in io.Reader, out io.Writer) {
	trace, err := traceback.ParseTraceback(in)
	if err != nil {
		panic(err)
	}
	stacks := trace.Stacks
	stacks = traceback.ExcludeGotest(stacks)

	if *filterGOROOT != "none" {
		stacks = traceback.ExcludeGoroot(stacks, *filterGOROOT == "chop")
	}
	if *filterTopOnly {
		stacks = traceback.ExcludeLowers(stacks)
	}
	switch *format {
	case "qfix":
		trace.Stacks = stacks
		traceback.Fprint(out, trace, traceback.PrintConfig{Quickfix: true})
	case "column":
	default:
		stacks = traceback.TrimSourcePrefix(stacks)
		trace.Stacks = stacks
		traceback.Fprint(out, trace, traceback.PrintConfig{})
	}
}
