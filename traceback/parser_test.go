package traceback

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

func ExampleParseStacks_go7725() {
	// stack trace example from go/issues/7725
	// http://code.google.com/p/go/issues/detail?id=7725
	// this stack trace contains several malformatted stacks such as runhome/XXX.
	// it may be caused by data handling misses or unknown bugs of golang's stack trace generator.
	printTraceSummary("testdata/go7725.txt")

	// Output:
	// Reason:SIGSEGV: segmentation violation
	// PC=0x4071dc
	// Goroutines:114 MinID:0 MaxID:13363
	// Status:sleep Count:100
	// Status:IO wait Count:3
	// Status:syscall Count:3
	// Status:chan receive Count:2
	// Status:select Count:2
	// Status:GC sweep wait Count:1
	// Status:finalizer wait Count:1
	// Status:garbage collection Count:1
	// Status:idle Count:1
	// Head:runtime.park Count:106
	// Head:runtime.notetsleepg Count:2
	// Head:runhome/dfc/go/src/pkg/runtime/time.goc:39 +0x31 fp=0x7f42b0033f70 Count:1
	// Head:runtime.gc Count:1
	// Head:runtime.park(0x413200, 0x12e22e0, 07f42b002ff58 Count:1
	// Head:runtime.park(0x413200, 0x12e22e0, org/v2/mgo/server.go:272 +0x110 fp=0x7f42b003cf98 Count:1
	// Head:scanblock Count:1
	// Head:syscall.Syscall Count:1
}
