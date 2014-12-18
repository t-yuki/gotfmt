package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/t-yuki/gotfmt/cmd"
	"github.com/t-yuki/gotfmt/webapp"
)

var flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

var (
	httpAddr = flags.String("http", "", "HTTP service address (e.g., ':6060')")
)

func init() {
	cmd.RegisterFlags(flags)
	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "gotfmt - Go Test formatter utility")
		fmt.Fprintf(os.Stderr, "Usage of %s [test|FILE]:\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "  If test is provided, it runs `go test` internally.")
		fmt.Fprintln(os.Stderr, "  Other flags or arguments will be passed to `go` command.")
		fmt.Fprintln(os.Stderr, "  If FILE is exists, it reads the test result from FILE.")
		fmt.Fprintln(os.Stderr)
		flags.PrintDefaults()
	}
}

func main() {
	goArgs := cmd.ParseFlags(flags, os.Args[1:])
	if *httpAddr != "" {
		webapp.ListenAndServe(*httpAddr)
		return
	}
	cmd.Main(goArgs)
}
