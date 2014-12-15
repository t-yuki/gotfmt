package traceback

import (
	"bytes"
	"fmt"
	"os"
	"sort"
)

func printTrace(filename string) (ignored *bytes.Buffer) {
	ignored = &bytes.Buffer{}
	data, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer data.Close()
	trace, err := ParseTraceback(data, ignored)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Reason:%s\n", trace.Reason)
	for _, s := range trace.Races {
		fmt.Printf("Race ID:%d Status:%s Calls:%d", s.ID, s.Status, len(s.Calls))
		if len(s.Calls) >= 1 {
			fmt.Printf(" Head:%s Line:%d", s.Calls[0].Func, s.Calls[0].Line)
		}
		fmt.Println()
	}
	for _, s := range trace.Stacks {
		fmt.Printf("ID:%d Status:%s Calls:%d", s.ID, s.Status, len(s.Calls))
		if len(s.Calls) >= 1 {
			fmt.Printf(" Head:%s Line:%d", s.Calls[0].Func, s.Calls[0].Line)
		}
		fmt.Println()
	}
	return
}

func printTraceSummary(filename string) (ignored *bytes.Buffer) {
	ignored = &bytes.Buffer{}
	defer func() {
		if e := recover(); e != nil {
			fmt.Fprintf(os.Stderr, "%s", ignored.String())
			panic(e)
		}
	}()
	data, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer data.Close()
	trace, err := ParseTraceback(data, ignored)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Reason:%s\n", trace.Reason)

	minID := 0xFFFF
	maxID := 0
	statusCounts := keyedCounters{}
	headFuncCounts := keyedCounters{}
	for _, s := range trace.Stacks {
		if minID > s.ID {
			minID = s.ID
		}
		if maxID < s.ID {
			maxID = s.ID
		}
		statusCounts.Add(string(s.Status), 1)
		if len(s.Calls) >= 1 {
			headFuncCounts.Add(s.Calls[0].Func, 1)
		}
	}
	fmt.Printf("Goroutines:%d MinID:%d MaxID:%d\n",
		len(trace.Stacks), minID, maxID)

	sort.Sort(sort.Reverse(statusCounts))
	for _, e := range statusCounts {
		fmt.Printf("Status:%s Count:%d\n", e.key, e.count)
	}

	sort.Sort(sort.Reverse(headFuncCounts))
	for _, e := range headFuncCounts {
		fmt.Printf("Head:%s Count:%d\n", e.key, e.count)
	}
	return
}

type keyedCounters []struct {
	key   string
	count int
}

func (a keyedCounters) Len() int {
	return len(a)
}

func (a *keyedCounters) Add(key string, delta int) {
	for i := range *a {
		if (*a)[i].key == key {
			(*a)[i].count += delta
			return
		}
	}
	*a = append(*a, struct {
		key   string
		count int
	}{key, delta})
}

func (a keyedCounters) Less(i, j int) bool {
	if a[i].count != a[j].count {
		return a[i].count < a[j].count
	}
	return a[i].key > a[j].key
}

func (a keyedCounters) Swap(i, j int) {
	t := a[i]
	a[i] = a[j]
	a[j] = t
}
