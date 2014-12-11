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

func Main(args []string) {
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
		err := Run(args)
		if err != nil {
			os.Exit(1)
		}
	}
}

func Run(args []string) (err error) {
	log := &bytes.Buffer{}
	wr := io.MultiWriter(os.Stdout, log)

	var in io.Reader
	var cmd *exec.Cmd
	if len(args) != 0 && args[0] != "test" {
		if _, err := os.Stat(args[0]); err == nil {
			if f, err := os.Open(args[0]); err == nil {
				in = f
				defer f.Close()
			}
		}
	}
	if len(args) != 0 && in == nil {
		cmd = exec.Command("go", args...)
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
	if in == nil {
		in = io.Reader(os.Stdin)
	}
	goErr := &bytes.Buffer{}
	in = io.TeeReader(in, goErr)

	tr := Convert(in, wr)

	if cmd != nil {
		err = cmd.Wait()
	}
	if tr != nil {
		_, h, isTTY := getScreenSize()
		logLines := countLines(log)
		if isTTY && h-3 < logLines {
			runPager(log)
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
