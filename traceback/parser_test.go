package traceback

import (
	"fmt"
	"os"
)

func ExampleParseStacks_deadlock() {
	printTrace("testdata/deadlock.txt")
	// Output:
	// Reason:fatal error: all goroutines are asleep - deadlock!
	// ID:1 Status:chan receive Calls:3 Head:testing.RunTests
	// ID:3 Status:chan receive (nil chan) Calls:3 Head:github.com/t-yuki/gotracetools/traceback/testdata.TestDeadlock
	// ID:4 Status:chan send (nil chan) Calls:2 Head:github.com/t-yuki/gotracetools/traceback/testdata.func·001
}

func ExampleParseStacks_timeout() {
	printTrace("testdata/timeout.txt")
	// Output:
	// Reason:panic: test timed out after 1s
	// ID:5 Status:running Calls:3 Head:runtime.panic
	// ID:1 Status:chan receive Calls:3 Head:testing.RunTests
	// ID:4 Status:sleep Calls:4 Head:time.Sleep
}

func ExampleParseStacks_timeout_early() {
	printTrace("testdata/timeout_early.txt")
	// Output:
	// Reason:panic: test timed out after 1us
	// ID:4 Status:running Calls:3 Head:runtime.panic
	// ID:1 Status:runnable Calls:10 Head:syscall.Syscall
}

func ExampleParseStacks_sigabrt() {
	printTrace("testdata/sigabrt.txt")
	// Output:
	// Reason:SIGABRT: abort
	// PC=0x424dd1
	// ID:1 Status:chan receive Calls:3 Head:testing.RunTests
	// ID:4 Status:chan receive Calls:3 Head:github.com/t-yuki/gotracetools/traceback/testdata.TestSIGABRT
}

func printTrace(filename string) {
	data, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer data.Close()
	trace, err := ParseTraceback(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Reason:%s\n", trace.Reason)
	for _, s := range trace.Stacks {
		fmt.Printf("ID:%d Status:%s Calls:%d", s.ID, s.Status, len(s.Calls))
		if len(s.Calls) >= 1 {
			fmt.Printf(" Head:%s", s.Calls[0].Func)
		}
		fmt.Println()
	}
}
