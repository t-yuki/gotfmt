package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
)

func Main(goArgs []string) {
	repeat, procs := *fNrepeat, *fProcs
	if *fNP != 0 {
		repeat = *fNP
	}
	if procs != 0 {
		os.Setenv("GOMAXPROCS", strconv.Itoa(procs))
	}

	// ignore sigquit
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, syscall.SIGQUIT)

	for i := 0; i < repeat; i++ {
		if *fNP != 0 {
			os.Setenv("GOMAXPROCS", strconv.Itoa(i+1))
		}
		err := Run(goArgs)
		if err != nil {
			os.Exit(1)
		}
	}
}

func Run(goArgs []string) (err error) {
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
			return err
		}
	}
	goErr := &bytes.Buffer{}
	in = io.TeeReader(in, goErr)

	Convert(in, os.Stdout)
	if cmd != nil {
		err = cmd.Wait()
	}
	return err
}
