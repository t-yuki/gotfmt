package traceback

import (
	"go/build"
	"strings"
)

// TrimSourcePrefix trims file path prefix from src.Calls.Source field.
func TrimSourcePrefix(srcs []Stack) []Stack {
	dst := make([]Stack, 0, len(srcs))
	for _, src := range srcs {
		s := src
		for i := range s.Calls {
			c := &s.Calls[i]
			for _, prefix := range []string{"src/pkg/", "src/"} {
				idx := strings.Index(c.Source, prefix)
				if idx != -1 {
					c.Source = c.Source[idx+len(prefix):]
				}
			}
		}
		dst = append(dst, s)
	}
	return dst
}

// ExcludeGotest filters out no meaningful stacks and calls caused by go test command and testing package.
func ExcludeGotest(srcs []Stack) (dst []Stack) {
	dst = filterGotestStacks(srcs)
	dst = filterGotestCalls(dst)
	return dst
}

// filterGotestStacks filters out no meaningful gotest stack.
func filterGotestStacks(src []Stack) []Stack {
	dst := make([]Stack, 0, len(src))
	for _, s := range src {
		if len(s.Calls) < 2 {
			dst = append(dst, s)
			continue
		}
		bottom := s.Calls[len(s.Calls)-1]
		// if this stack is testmain, it isn't meaningful so filter out.
		if strings.Contains(bottom.Source, "_test/_testmain.go") {
			continue
		}
		// if this stack is testing timer, it isn't meaningful so filter out.
		if strings.Contains(bottom.Source, "time/sleep.go") {
			if strings.Contains(s.Calls[len(s.Calls)-2].Source, "testing/testing.go") {
				continue
			}
		}
		dst = append(dst, s)
	}
	return dst
}

func filterGotestCalls(srcs []Stack) []Stack {
	dst := make([]Stack, 0, len(srcs))
	for _, src := range srcs {
		s := src
		for i := len(src.Calls) - 1; i >= 0; i-- {
			if strings.Contains(src.Calls[i].Source, "testing/testing.go") {
				continue
			}
			s.Calls = src.Calls[0 : i+1]
			break
		}
		dst = append(dst, s)
	}
	return dst
}

// ExcludeGoroot filters out GOROOT function calls form the top of stacks.
// If `preserveOne` is set, it leaves the nearest GOROOT function call.
func ExcludeGoroot(srcs []Stack, preserveOne bool) []Stack {
	dst := make([]Stack, 0, len(srcs))
	for _, src := range srcs {
		s := src
		s.Calls = nil
		for i := range src.Calls {
			c := &src.Calls[i]
			if strings.HasPrefix(c.Source, build.Default.GOROOT+"/") {
				continue
			}
			// a special value for official binary build
			if strings.HasPrefix(c.Source, "/usr/local/go/src/pkg/") {
				continue
			}
			if i != 0 {
				if preserveOne {
					s.Calls = src.Calls[i-1:]
				} else {
					s.Calls = src.Calls[i:]
				}
			} else {
				s.Calls = src.Calls
			}
			break
		}
		if s.Calls != nil {
			dst = append(dst, s)
		}
	}
	return dst
}

// ExcludeLowers filters calls of stacks by excluding function calls that is not the top of stack.
func ExcludeLowers(srcs []Stack) []Stack {
	dst := make([]Stack, 0, len(srcs))
	for _, src := range srcs {
		s := src
		s.Calls = []Call{src.Calls[0]}
		dst = append(dst, s)
	}
	return dst
}
