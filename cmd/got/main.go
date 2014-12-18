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
	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "got - Go Test runner utility")
		fmt.Fprintf(os.Stderr, "Usage of %s [FILE]:\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "  If FILE is exists, it does not run `go test`.")
		fmt.Fprintln(os.Stderr, "  Instead, it reads the test result from FILE.")
		fmt.Fprintln(os.Stderr, "  Any other flags or arguments will be passed to `go test` command.")
		fmt.Fprintln(os.Stderr)
		flags.PrintDefaults()
	}
}

func main() {
	mode := "test"
	args := cmd.ParseFlags(flags, os.Args[1:])

	if len(args) != 0 {
		if _, err := os.Stat(args[0]); err == nil {
			mode = args[0]
			args = args[1:]
		}
	}
	cmd.Main(append([]string{mode}, args...))
}
