// This is a derivative work of roger peppe's stackgraph command.
// For more details, see http://code.google.com/p/rog-go/
//

package traceback

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// ParseTraceback parse stacktrace of go executables.
func ParseTraceback(r io.Reader) (*Traceback, error) {
	trace := &Traceback{}
	s := bufio.NewScanner(r)

	stacks := make([]Stack, 0, 100)
	re1 := regexp.MustCompile(`(panic|fatal error|SIG\w+):`)
	re2 := regexp.MustCompile(`goroutine (\d+) \[([\w ]+)\]:`)
	for s.Scan() {
		line := s.Text()
		if trace.Reason == "" {
			strs := re1.FindStringSubmatch(line)
			if strs != nil && strs[0] != "" && strs[1] != "" {
				trace.Reason = line
				for s.Scan() {
					line := s.Text()
					if line == "" {
						// empty line signifies end of a stack
						break
					}
					trace.Reason += "\n" + line
				}
				continue
			}
		}

		strs := re2.FindStringSubmatch(line)
		if strs == nil || strs[0] == "" || strs[1] == "" || strs[2] == "" {
			continue
		}
		stack := Stack{}
		stack.ID, _ = strconv.Atoi(strs[1])
		stack.Status = StackStatus(strs[2])
		for s.Scan() {
			line := s.Text()
			if line == "" {
				// empty line signifies end of a stack
				break
			}
			if strings.HasPrefix(line, "exit status") {
				break
			}
			if strings.Contains(line, "  ") {
				// Looks like a register dump.
				// TODO better heuristic here.
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
			if !s.Scan() {
				break
			}
			line = s.Text()
			if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, "        ") {
				line = strings.TrimPrefix(line, "        ")
				line = strings.TrimPrefix(line, "\t")
				if i := strings.LastIndex(line, ":"); i >= 0 {
					call.Source = line[0:i]
					line = line[i:]
				}
				if i := strings.LastIndex(line, " +"); i >= 1 {
					call.Line, _ = strconv.Atoi(line[1:i])
				}
			}
			stack.Calls = append(stack.Calls, call)
		}
		if len(stack.Calls) > 0 {
			stacks = append(stacks, stack)
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
