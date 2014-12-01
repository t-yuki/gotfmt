package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
var fNrepeat = flags.Int("n", 1, "repeat the test N times while it passes")
var fProcs = flags.Int("p", 0, "set GOMAXPROCS")
var fNP = flags.Int("np", 0, "similar to a combination of `-n` and `-p` but increment GOMAXPROCS from 1 for each repeat")

func init() {
	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "got - Go Test runner utility")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flags.PrintDefaults()
		fmt.Fprintln(os.Stderr, "Any other flags will be passed to `go test` command or `gotb` command.")
	}
}

type got struct {
	test *exec.Cmd
	gotb *exec.Cmd
}

func main() {
	args := make([]string, 0, 10)
	testArgs := make([]string, 0, 10)
	gotbArgs := make([]string, 0, 10)

	var skip bool
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if skip {
			testArgs = append(testArgs, arg)
			continue
		}
		pair := strings.SplitN(arg, "=", 2)
		key := strings.TrimLeft(pair[0], "-")
		name := strings.TrimRight(key, "0123456789")
		value := strings.TrimPrefix(key, name)
		if value != "" && len(pair) == 1 {
			arg = "-" + name + "=" + value
			pair = []string{name, value}
		}
		switch {
		case strings.HasPrefix(name, "filter.") || strings.HasPrefix(name, "out."):
			gotbArgs = append(gotbArgs, arg)
		case flags.Lookup(name) != nil || name == "h" || name == "help":
			args = append(args, arg)
			if i+1 < len(os.Args) && len(pair) == 1 {
				fv, ok := flags.Lookup(name).Value.(interface {
					IsBoolFlag() bool
				})
				if !ok || !fv.IsBoolFlag() {
					args = append(args, os.Args[i+1])
					i++
				}
			}
		case name == "":
			skip = true
		default:
			testArgs = append(testArgs, arg)
		}
	}
	flags.Parse(args)

	repeat, procs := *fNrepeat, *fProcs
	if *fNP != 0 {
		repeat = *fNP
	}
	if procs != 0 {
		os.Setenv("GOMAXPROCS", strconv.Itoa(procs))
	}

	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, syscall.SIGQUIT)

	for i := 0; i < repeat; i++ {
		if *fNP != 0 {
			os.Setenv("GOMAXPROCS", strconv.Itoa(i+1))
		}
		var err error
		testend := make(chan error, 1)
		var g got
		g.start(testend, testArgs, gotbArgs)
		select {
		case <-sigquit:
			err = <-testend
		case err = <-testend:
			break
		}
		g.stop()
		if err != nil {
			break
		}
	}
}

func (g *got) start(endch chan<- error, gotestargs []string, gotbargs []string) {
	goargs := []string{"test"}
	goargs = append(goargs, gotestargs...)
	g.test = exec.Command("go", goargs...)
	g.test.Stdin = os.Stdin
	g.test.Stdout = os.Stdout
	terr, err := g.test.StderrPipe()
	if err != nil {
		panic(err)
	}

	g.gotb = exec.Command("gotb", gotbargs...)
	g.gotb.Stderr = os.Stderr
	g.gotb.Stdout = os.Stdout
	tbin, err := g.gotb.StdinPipe()
	if err != nil {
		panic(err)
	}

	go io.Copy(tbin, terr)
	g.test.Start()
	g.gotb.Start()

	go func() {
		err := g.test.Wait()
		tbin.Close()
		endch <- err
	}()
}

func (g *got) stop() {
	g.gotb.Wait()
}
