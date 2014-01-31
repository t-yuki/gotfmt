package traceback

import (
	"fmt"
	"os"
)

func ExampleParseStacks_data1() {
	data, err := os.Open("testdata/data1.txt")
	if err != nil {
		panic(err)
	}
	trace, err := ParseTraceback(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", trace.Reason)
	for _, s := range trace.Stacks {
		fmt.Printf("%d %s\n", s.ID, s.Status)
		for _, c := range s.Calls {
			fmt.Printf("    %v\n", c)
		}
	}
	// Output:
	// fatal error: all goroutines are asleep - deadlock!
	// 1 chan receive
	//     {testing.RunTests /usr/local/go/src/pkg/testing/testing.go 472 [5607464 6582432 1 1 1]}
	//     {testing.Main /usr/local/go/src/pkg/testing/testing.go 403 [5607464 6582432 1 1 6615616]}
	//     {main.main github.com/t-yuki/mygosandbox/go2qfix/_test/_testmain.go 47 []}
	// 3 semacquire
	//     {sync.runtime_Syncsemacquire /usr/local/go/src/pkg/runtime/sema.goc 257 [833492435008]}
	//     {sync.(*Cond).Wait /usr/local/go/src/pkg/sync/cond.go 62 [833492434992]}
	//     {io.(*pipe).read /usr/local/go/src/pkg/io/pipe.go 52 [833492434944 833492439040 4096 4096 0]}
	//     {io.(*PipeReader).Read /usr/local/go/src/pkg/io/pipe.go 134 [833492091088 833492439040 4096 4096 4317609]}
	//     {bufio.(*Scanner).Scan /usr/local/go/src/pkg/bufio/scan.go 165 [833492410592 4096]}
	//     {github.com/t-yuki/mygosandbox/go2qfix.TestConvertEmpty /home/sey/gopath/src/github.com/t-yuki/mygosandbox/go2qfix/go2qfix_test.go 16 [833492430848]}
	//     {testing.tRunner /usr/local/go/src/pkg/testing/testing.go 391 [833492430848 6582432]}
	//     {testing.RunTests /usr/local/go/src/pkg/testing/testing.go 471 []}
}

func ExampleParseStacks_data2() {
	data, err := os.Open("testdata/data2.txt")
	if err != nil {
		panic(err)
	}
	trace, err := ParseTraceback(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", trace.Reason)
	for _, s := range trace.Stacks {
		fmt.Printf("%d %s\n", s.ID, s.Status)
		for _, c := range s.Calls {
			fmt.Printf("    %v\n", c)
		}
	}
	// Output:
	// panic: test timed out after 1us
	// 5 running
	//     {runtime.panic /usr/local/go/src/pkg/runtime/panic.c 266 [5032096 833492317296]}
	//     {testing.func·007 /usr/local/go/src/pkg/testing/testing.go 596 []}
	//     {time.goFunc /usr/local/go/src/pkg/time/sleep.go 123 []}
	// 1 runnable
	//     {testing.func·004 /usr/local/go/src/pkg/testing/example.go 79 []}
	//     {testing.runExample /usr/local/go/src/pkg/testing/example.go 100 [5436144 24 5594232 5631280 1074]}
	//     {testing.RunExamples /usr/local/go/src/pkg/testing/example.go 36 [5594312 6571488 1 1 1]}
	//     {testing.Main /usr/local/go/src/pkg/testing/testing.go 404 [5594312 6603328 0 0 6603328]}
	//     {main.main github.com/t-yuki/gotracetools/traceback/_test/_testmain.go 47 []}
}
