package main

import (
	"flag"
	"github.com/t-yuki/gotracetools/traceback"
	"io"
	"os"
	"os/signal"
	"syscall"
)

var (
	excludeGOROOT = flag.Bool("R", false, "exclude GOROOT functions completely")
	includeGOROOT = flag.Bool("r", false, "include GOROOT functions")
	topOnly       = flag.Bool("t", false, "print top functions only (implies `-R` when `r` == false)")
	quickfix      = flag.Bool("q", false, "print with vim quickfix list format (implies `-t -R`). Hint: gotb -q | vim - -c :cb! -c :copen")
)

func main() {
	flag.Parse()
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, syscall.SIGQUIT)
	convert(os.Stdin, os.Stdout)
}

func convert(in io.Reader, out io.WriteCloser) {
	if *quickfix {
		*topOnly = true
	}
	if *topOnly && !*includeGOROOT {
		*excludeGOROOT = true
	}

	trace, err := traceback.ParseTraceback(in)
	if err != nil {
		panic(err)
	}
	stacks := trace.Stacks
	stacks = traceback.ExcludeGotest(stacks)
	if !*includeGOROOT {
		stacks = traceback.ExcludeGoroot(stacks, !*excludeGOROOT)
	}
	if *topOnly {
		stacks = traceback.ExcludeLowers(stacks)
	}
	if *quickfix {
		trace.Stacks = stacks
		traceback.Fprint(out, trace, traceback.PrintConfig{Quickfix: true})
	} else {
		stacks = traceback.TrimSourcePrefix(stacks)
		trace.Stacks = stacks
		traceback.Fprint(out, trace, traceback.PrintConfig{})
	}
	out.Close()
}
