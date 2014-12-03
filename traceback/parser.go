// This is a derivative work of roger peppe's stackgraph command.
// For more details, see http://code.google.com/p/rog-go/
//

package traceback

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"regexp/syntax"
	"strconv"
	"strings"
)

// ParseTraceback parses the per-goroutine stacktraces generated by GOTRACEBACK=1 environment.
// It reads `r` and constructs a `Traceback` struct to return.
// It also writes unrelated lines to `w`.
func ParseTraceback(r io.Reader, w io.Writer) (*Traceback, error) {
	trace := &Traceback{}
	s := bufio.NewScanner(r)
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	stacks := make([]Stack, 0, 100)
	for s.Scan() {
		line := s.Text()
		if trace.Reason == "" {
			if reason, ok := parseReason(s, bw, line); ok {
				trace.Reason = reason
			}
			// no reason???
			continue
		}

		if stack, ok := parseStack(s, bw, line); ok {
			if len(stack.Calls) > 0 {
				stacks = append(stacks, stack)
			}
			continue
		}
	}
	trace.Stacks = stacks
	return trace, nil
}

func parseArgs(argList string) []uint64 {
	argList = strings.TrimSuffix(argList, ", ...")
	if argList == "" {
		return nil
	}
	parts := strings.Split(argList, ", ")
	args := make([]uint64, len(parts))
	for i, a := range parts {
		n, err := strconv.ParseUint(a, 0, 64)
		if err != nil {
			panic(fmt.Errorf("cannot parse %q (from %q)", a, argList))
			n = 0xdeadbeef
		}
		args[i] = n
	}
	return args
}

// _+FAIL: is user-defined reason so that traceback is not generated by go runtime
var reasonRE = regexp.MustCompile(`(panic|fatal error|SIG\w+|_+(\w+_+)?FAIL):`)

func parseReason(s *bufio.Scanner, w *bufio.Writer, line string) (reason string, ok bool) {
	strs := reasonRE.FindStringSubmatch(line)
	if strs == nil || strs[0] == "" || strs[1] == "" {
		w.WriteString(line + "\n")
		return "", false
	}
	reason = strings.TrimLeft(line, "_")
	for s.Scan() {
		line := s.Text()
		if line == "" {
			// empty line signifies end of a stack
			break
		}
		reason += "\n" + line
	}
	return reason, true
}

var stackRE = regexp.MustCompile(`goroutine (\d+) \[([\w\d, ()]+)\]:`)

func parseStack(s *bufio.Scanner, w *bufio.Writer, line string) (stack Stack, ok bool) {
	strs := stackRE.FindStringSubmatch(line)
	if strs == nil || strs[0] == "" || strs[1] == "" || strs[2] == "" {
		w.WriteString(line + "\n")
		return stack, false
	}
	stack.ID, _ = strconv.Atoi(strs[1])
	stack.Status = StackStatus(strs[2])
	for s.Scan() {
		line := s.Text()
		if endOfTraceback(line) {
			break
		}
		if shouldIgnore(line) {
			w.WriteString(line + "\n")
			continue
		}
		var call Call
		if strings.HasSuffix(line, ")") {
			if i := strings.LastIndex(line, "("); i > 0 {
				call.Args = parseArgs(line[i+1 : len(line)-1])
				line = line[0:i]
			}
		}
		call.Func = strings.TrimPrefix(line, "created by ")
		for {
			if !s.Scan() {
				return stack, true
			}
			line = s.Text()
			if !shouldIgnore(line) {
				break
			}
			w.WriteString(line + "\n")
		}
		if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, "        ") {
			line = strings.TrimPrefix(line, "        ")
			line = strings.TrimPrefix(line, "\t")
			if i := strings.LastIndex(line, ":"); i >= 0 {
				call.Source = line[0:i]
				line = line[i:]
			}
			if i := strings.LastIndex(line, " +"); i >= 1 {
				call.Line, _ = strconv.Atoi(line[1:i])
			} else {
				call.Line, _ = strconv.Atoi(line[1:])
			}
		}
		stack.Calls = append(stack.Calls, call)
	}
	return stack, true
}

func endOfTraceback(line string) bool {
	if line == "----- stack segment boundary -----" {
		return false
	}
	if line == "" {
		// empty line signifies end of a stack
		return true
	}
	if strings.HasPrefix(line, "exit status") {
		return true
	}
	if strings.HasPrefix(line, "FAIL") {
		return true
	}
	if strings.IndexAny(line, "=-?") == 0 {
		return true
	}
	return false
}

func shouldIgnore(line string) bool {
	if line == "" {
		return true
	}
	if syntax.IsWordChar(rune(line[0])) && strings.Contains(line, "  0x") {
		// Looks like a register dump.
		// TODO better heuristic here.
		return true
	}
	if str := strings.TrimSpace(line); !syntax.IsWordChar(rune(str[0])) && str[0] != '/' {
		// Looks like a mixed line.
		return true
	}
	return false
}
