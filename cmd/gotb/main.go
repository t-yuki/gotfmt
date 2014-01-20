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
	topOnly       = flag.Bool("t", false, "show top functions only (implies `-R` when `r` == false)")
)

func main() {
	flag.Parse()
	convert(os.Stdin, os.Stdout)
}

func convert(in io.Reader, out io.WriteCloser) {
	stacks, _ := traceback.ParseStacks(in)
	stacks = traceback.ExcludeGotest(stacks)
	if !*includeGOROOT {
		stacks = traceback.ExcludeGoroot(stacks, !*excludeGOROOT && !*topOnly)
	}
	if *topOnly {
		stacks = traceback.ExcludeLowers(stacks)
	}
	stacks = traceback.TrimSourcePrefix(stacks)
	traceback.Fprint(out, stacks)
	out.Close()
}
