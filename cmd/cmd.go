package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func Main(args []string) {
	repeat, procs := *fNrepeat, *fProcs
	if *fNP != 0 {
		repeat = *fNP
	}
	if procs != 0 {
		os.Setenv("GOMAXPROCS", strconv.Itoa(procs))
	}

	testbin := ""

	if repeat > 1 {
		file, err := ioutil.TempFile("", "gotfmt")
		if err != nil {
			panic(err)
		}
		defer os.Remove(file.Name())
		err = build(file.Name(), args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		testbin = file.Name()
	}

	for i := 0; i < repeat; i++ {
		if *fNP != 0 {
			os.Setenv("GOMAXPROCS", strconv.Itoa(i+1))
		}
		err := Run(args, testbin)
		if err != nil {
			os.Exit(1)
		}
	}
}

func Run(args []string, testbin string) (err error) {
	log := &bytes.Buffer{}
	logger := io.MultiWriter(os.Stdout, log)

	in := io.Reader(os.Stdin)
	testErrCh := (chan error)(nil)
	if len(args) != 0 {
		if args[0] == "test" { // run gotest, use prebuilt if available
			in, testErrCh, err = runGotest(testbin, args, logger)
			if err != nil {
				return err
			}
		} else if _, err := os.Stat(args[0]); err == nil { // read stacktrace from a file
			if f, err := os.Open(args[0]); err == nil {
				in = f
				defer f.Close()
			}
		}
	} // otherwise, use stdin

	goErr := &bytes.Buffer{}
	in = io.TeeReader(in, goErr)

	trace := Convert(in, logger)

	if testErrCh != nil {
		testErrCh <- nil
		err = <-testErrCh
	}
	if trace != nil {
		_, h, isTTY := getScreenSize()
		logLines := countLines(log)
		if isTTY && h-3 < logLines {
			runPager(log)
		}
	}
	return err
}

func appendOptPrefix(args []string) []string {
	ret := make([]string, 0, len(args))
	for _, arg := range args {
		pair := strings.Split(arg, "=")
		// TODO: skip next arg if len(pair) == 0 and pair[0] is not bool type
		switch pair[0] {
		case "-race", "-cover", "-covermode", "-coverpkg":
			continue // ignore build time flag
		case "-bench", "-benchmem", "-benchtime":
		case "-blockprofile", "-blockprofilerate":
		case "-coverprofile":
		case "-cpu", "-cpuprofile":
		case "-memprofile", "-memprofilerate":
		case "-outputdir", "-parallel", "-run", "-short", "-timeout", "-v":
			break // append test. prefix
		default:
			ret = append(ret, arg)
			continue
		}
		ret = append(ret, "-test."+arg[1:])
	}
	return ret
}

func runGotest(testbin string, args []string, logger io.Writer) (errlogOut io.Reader, execResult chan error, startError error) {
	prebuilt := testbin != ""
	if prebuilt {
		if args[0] == "test" {
			args = args[1:]
		}
		args = appendOptPrefix(args)
	} else {
		testbin = "go"
	}

	// ignore sigquit, sigint and sigterm during gotest
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGQUIT)
	signal.Notify(sigch, os.Interrupt)

	cmd := exec.Command(testbin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = logger
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	resultChan := make(chan error, 0)
	go catchCommandResult(cmd, resultChan, sigch, prebuilt)
	return stderr, resultChan, nil
}

func catchCommandResult(cmd *exec.Cmd, resultChan chan error, sigch chan os.Signal, prebuilt bool) {
	<-resultChan // wait for Convert is done
	err := cmd.Wait()
	signal.Stop(sigch)
	if prebuilt {
		if err != nil {
			exiterr := err.(*exec.ExitError)
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				fmt.Fprintf(os.Stderr, "exit status %d\n", status.ExitStatus())
			}
			fmt.Fprintf(os.Stderr, "FAIL\n") // TODO: write package name and time or fail reason
		} else {
			fmt.Fprintf(os.Stderr, "ok\n") // TODO: write package name and time
		}
	}
	resultChan <- err
}

func getScreenSize() (w, h int, ok bool) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, false
	}
	fmt.Sscanf(string(out), "%d %d\n", &h, &w)
	if w == 0 || h == 0 {
		return w, h, false
	}
	return w, h, true
}

func countLines(buf *bytes.Buffer) (n int) {
	for _, v := range buf.Bytes() {
		if v == '\n' {
			n++
		}
	}
	return
}

func runPager(buf *bytes.Buffer) {
	args := []string{"-R"}
	line, ok := findFirstGoroutineLine(buf)
	if ok {
		if line > 20 {
			line -= 20
		} else {
			line = 0
		}
		args = append(args, "+"+strconv.Itoa(line))
	}

	cmd := exec.Command("less", args...)
	cmd.Stdin = buf
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func findFirstGoroutineLine(buf *bytes.Buffer) (n int, ok bool) {
	buf = bytes.NewBuffer(buf.Bytes())
	scan := bufio.NewScanner(buf)
	for scan.Scan() {
		n++
		if strings.HasPrefix(scan.Text(), "goroutine ") {
			return n, true
		}
	}
	return 0, false
}
