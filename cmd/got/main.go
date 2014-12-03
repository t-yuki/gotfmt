package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/t-yuki/gotfmt/cmd"
)

var flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

func init() {
	cmd.RegisterFlags(flags)
	t := flags.Lookup("t")
	t.DefValue = "text"
	t.Value.Set("text")
	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "got - Go Test runner utility")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flags.PrintDefaults()
		fmt.Fprintln(os.Stderr, "Any other flags or arguments will be passed to `go test` command.")
	}
}

func main() {
	testArgs := cmd.ParseFlags(flags)
	cmd.Main(append([]string{"test"}, testArgs...))
}
