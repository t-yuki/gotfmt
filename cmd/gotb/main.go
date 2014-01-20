package main

import (
	"flag"
	"github.com/t-yuki/gotracetools/traceback"
	"io"
	"os"
)

var (
	excludeGOROOT = flag.Bool("R", false, "exclude GOROOT functions completely")
	includeGOROOT = flag.Bool("r", false, "include GOROOT functions")
	topOnly       = flag.Bool("t", false, "print top functions only (implies `-R` when `r` == false)")
	quickfix      = flag.Bool("q", false, "print with vim quickfix list format (implies `-t -R`). Hint: gotb -q | vim - -c :cb!")
)

func main() {
	flag.Parse()
	convert(os.Stdin, os.Stdout)
}

func convert(in io.Reader, out io.WriteCloser) {
	if *quickfix {
		*topOnly = true
	}
	if *topOnly && !*includeGOROOT {
		*excludeGOROOT = true
	}

	stacks, _ := traceback.ParseStacks(in)
	stacks = traceback.ExcludeGotest(stacks)
	if !*includeGOROOT {
		stacks = traceback.ExcludeGoroot(stacks, !*excludeGOROOT)
	}
	if *topOnly {
		stacks = traceback.ExcludeLowers(stacks)
	}
	if *quickfix {
		traceback.Fprint(out, stacks, traceback.PrintConfig{Quickfix: true})
	} else {
		stacks = traceback.TrimSourcePrefix(stacks)
		traceback.Fprint(out, stacks, traceback.PrintConfig{})
	}
	out.Close()
}
