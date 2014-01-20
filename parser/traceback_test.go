package parser

import (
	"fmt"
	"os"
)

func ExampleParseStacks_data1() {
	data1, err := os.Open("testdata/data1.txt")
	if err != nil {
		panic(err)
	}
	stacks, err := ParseStacks(data1)
	if err != nil {
		panic(err)
	}
	for _, s := range stacks {
		fmt.Printf("%d %s\n", s.ID, s.Status)
		for _, c := range s.Calls {
			fmt.Printf("    %v\n", c)
		}
	}
	// Output:
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
