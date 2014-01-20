package main

import (
	"flag"
	"github.com/t-yuki/gotracetools/traceback"
	"io"
	"os"
)

func main() {
	flag.String("r", "", "include GOROOT functions")
	flag.Parse()
	convert(os.Stdin, os.Stdout)
}

func convert(in io.Reader, out io.WriteCloser) {
	stacks, _ := traceback.ParseStacks(in)
	stacks = traceback.ExcludeGotest(stacks)
	stacks = traceback.ExcludeGoroot(stacks)
	stacks = traceback.TrimSourcePrefix(stacks)
	traceback.Fprint(out, stacks)
	out.Close()
}
