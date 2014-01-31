package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

type got struct {
	test *exec.Cmd
	gotb *exec.Cmd
}

func main() {
	gotestargs := make([]string, 0, 10)
	gotbargs := make([]string, 0, 10)

	var skip bool
	for _, s := range os.Args[1:] {
		if skip {
			gotestargs = append(gotestargs, s)
			continue
		}
		switch strings.SplitN(s, "=", 2)[0] {
		case "-R", "-r", "-t", "-q":
			gotbargs = append(gotbargs, s)
		case "--":
			skip = true
		default:
			gotestargs = append(gotestargs, s)
		}
	}

	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, syscall.SIGQUIT)

	testend := make(chan struct{}, 1)
	var g got
	g.start(testend, gotestargs, gotbargs)
	select {
	case <-sigquit:
		<-testend
	case <-testend:
		break
	}
	g.stop()
}

func (g *got) start(endch chan<- struct{}, gotestargs []string, gotbargs []string) {
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
		g.test.Wait()
		fmt.Fprintln(os.Stdout)
		tbin.Close()
		endch <- struct{}{}
	}()
}

func (g *got) stop() {
	g.gotb.Wait()
}
