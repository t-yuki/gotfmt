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
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flags.PrintDefaults()
		fmt.Fprintln(os.Stderr, "Any other flags or arguments will be passed to `go` command.")
	}
}

func main() {
	goArgs := cmd.ParseFlags(flags)
	if *httpAddr != "" {
		webapp.ListenAndServe(*httpAddr)
		return
	}
	cmd.Main(goArgs)
}
