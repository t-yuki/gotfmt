package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
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
	log := &bytes.Buffer{}
	wr := io.MultiWriter(os.Stdout, log)

	in := io.Reader(os.Stdin)
	var cmd *exec.Cmd
	if len(goArgs) != 0 {
		cmd = exec.Command("go", goArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = wr
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

	Convert(in, wr)

	if cmd != nil {
		err = cmd.Wait()
		if err != nil {
			_, h, ok := getScreenSize()
			logLines := countLines(log)
			if ok && h-3 < logLines {
				runPager(log)
			}
		}
	}
	return err
}

func getScreenSize() (w, h int, ok bool) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, false
	}
	fmt.Sscanf(string(out), "%d %d\n", &w, &h)
	if w == 0 || h == 0 {
		return w, h, false
	}
	return w, h, true
}

func countLines(buf *bytes.Buffer) (n int) {
	for i := range buf.Bytes() {
		if i == '\n' {
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
