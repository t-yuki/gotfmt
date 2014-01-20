package traceback

import (
	"strings"
)

// FilterGotest excludes gotest related stacks and calls.
func FilterGotest(src []*Stack) (dst []*Stack) {
	dst = filterGotestStacks(src)
	dst = filterGotestCalls(dst)
	return dst
}

func filterGotestStacks(src []*Stack) []*Stack {
	dst := make([]*Stack, 0, len(src))
	for _, s := range src {
		if len(s.Calls) < 2 {
			dst = append(dst, s)
			continue
		}
		last := s.Calls[len(s.Calls)-1]
		if strings.Contains(last.Source, "_test/_testmain.go") {
			continue
		}
		if strings.Contains(last.Source, "time/sleep.go") {
			if strings.Contains(s.Calls[len(s.Calls)-2].Source, "testing/testing.go") {
				continue
			}
		}
		dst = append(dst, s)
	}
	return dst
}

func filterGotestCalls(srcs []*Stack) []*Stack {
	dst := make([]*Stack, 0, len(srcs))
	for _, src := range srcs {
		s := *src
		for i := len(src.Calls) - 1; i >= 0; i-- {
			if strings.Contains(src.Calls[i].Source, "testing/testing.go") {
				continue
			}
			s.Calls = src.Calls[0 : i+1]
			break
		}
		dst = append(dst, &s)
	}
	return dst
}
