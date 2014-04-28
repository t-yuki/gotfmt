package main

import (
	"bytes"
	flagp "flag"
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
	excludeGOROOT = flag.Bool("R", false, "exclude GOROOT functions completely")
	includeGOROOT = flag.Bool("r", false, "include GOROOT functions")
	topOnly       = flag.Bool("t", false, "print top functions only (implies `-R` when `r` == false)")
	quickfix      = flag.Bool("q", false, "print with vim quickfix list format (implies `-t -R`). Hint: gotb -q | vim - -c :cb! -c :copen")
	httpAddr      = flag.String("http", "", "HTTP service address (e.g., ':6060')")
)

func main() {
	goArgs, gotfmtArgs := preprocessArgs(os.Args[1:])
	flag.Parse(gotfmtArgs)

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
		err := cmd.Wait()
		if err != nil {
			io.Copy(os.Stderr, goErr)
			os.Stderr.WriteString(err.Error() + "\n")
		}
	}
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
