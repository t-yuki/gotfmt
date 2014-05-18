package main

import (
	"bytes"
	"encoding/json"
	flagp "flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/t-yuki/gotfmt/traceback"
	"github.com/t-yuki/gotfmt/webapp"
)

var flag = flagp.NewFlagSet(os.Args[0], flagp.ExitOnError)

var (
	help   = flag.Bool("h", false, "show this help")
	format = flag.String("t", "text", `output format
	text: pretty formatted text format
	qfix: vim quickfix output format with errorformat: '%f:%l:\ %m'. you should use with 'nostd,notest,top' filters
	json: JSON format`)
	filter = flag.String("f", "", `stack trace filters by comma-separated list
	trimstd:  exclude GOROOT function calls but leave one
	nostd:    exclude GOROOT function calls completely
	notest:   exclude testing function calls
	top:      remove lower function calls`)
	httpAddr = flag.String("http", "", "HTTP service address (e.g., ':6060')")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	goArgs, gotfmtArgs := preprocessArgs(os.Args[1:])
	flag.Parse(gotfmtArgs)

	if *help {
		flag.Usage()
		return
	}

	if *httpAddr != "" {
		webapp.ListenAndServe(*httpAddr)
		return
	}

	// ignore sigquit
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, syscall.SIGQUIT)

	run(goArgs)
}

func preprocessArgs(args []string) (goArgs []string, gotfmtArgs []string) {
	goArgs = make([]string, 0, 10)
	gotfmtArgs = make([]string, 0, 10)
	for i, s := range args {
		if s == "--" || s == "=" {
			goArgs = append(goArgs, args[i:]...)
			return
		}

		name := strings.SplitN(s, "=", 2)[0]
		if flag.Lookup(name[1:]) != nil {
			gotfmtArgs = append(gotfmtArgs, s)
		} else {
			goArgs = append(goArgs, s)
		}
	}
	return
}

func run(goArgs []string) {
	in := io.Reader(os.Stdin)
	var cmd *exec.Cmd
	if len(goArgs) != 0 {
		cmd = exec.Command("go", goArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		stderr, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		in = stderr
		err = cmd.Start()
		if err != nil {
			return
		}
	}
	goErr := &bytes.Buffer{}
	in = io.TeeReader(in, goErr)

	convert(in, os.Stdout)
	if cmd != nil {
		_ = cmd.Wait()
	}
}

func convert(in io.Reader, out io.WriteCloser) {
	trace, err := traceback.ParseTraceback(in)
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
		stacks = traceback.TrimSourcePrefix(stacks)
		trace.Stacks = stacks
		traceback.Fprint(out, trace, traceback.PrintConfig{})
	}
	out.Close()
}
