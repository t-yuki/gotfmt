package traceback

import (
	"os"
	"testing"
)

func TestExcludesGotest1(t *testing.T) {
	data, err := os.Open("testdata/data1.txt")
	if err != nil {
		panic(err)
	}
	trace, err := ParseTraceback(data)
	if err != nil {
		panic(err)
	}
	stacks := ExcludeGotest(trace.Stacks)
	if len(stacks) != 1 {
		t.Error(stacks)
	}
	if stacks[0].ID != 3 {
		t.Error(stacks[0])
	}
	if len(stacks[0].Calls) != 6 {
		t.Log(stacks[0])
		t.Error(len(stacks[0].Calls))
	}
}

func TestExcludesGotest2(t *testing.T) {
	data, err := os.Open("testdata/data2.txt")
	if err != nil {
		panic(err)
	}
	trace, err := ParseTraceback(data)
	if err != nil {
		panic(err)
	}
	stacks := trace.Stacks
	stacks = ExcludeGotest(stacks)
	if len(stacks) != 0 {
		t.Error(stacks)
	}
}

func TestExcludesGoroot1(t *testing.T) {
	data, err := os.Open("testdata/data1.txt")
	if err != nil {
		panic(err)
	}
	trace, err := ParseTraceback(data)
	if err != nil {
		panic(err)
	}
	stacks := trace.Stacks
	stacks = ExcludeGotest(stacks)
	stacks = ExcludeGoroot(stacks, true)
	if len(stacks) != 1 {
		t.Error(stacks)
	}
	if stacks[0].ID != 3 {
		t.Error(stacks[0])
	}
	if len(stacks[0].Calls) != 2 {
		t.Log(stacks[0])
		t.Error(len(stacks[0].Calls))
	}
}

func TestExcludesGoroot3(t *testing.T) {
	data, err := os.Open("testdata/data3.txt")
	if err != nil {
		panic(err)
	}
	trace, err := ParseTraceback(data)
	if err != nil {
		panic(err)
	}
	stacks := trace.Stacks
	stacks = ExcludeGotest(stacks)
	stacks = ExcludeGoroot(stacks, false)
	if len(stacks) != 1 {
		t.Error(stacks)
	}
	if stacks[0].ID != 10 {
		t.Error(stacks[0])
	}
	if len(stacks[0].Calls) != 1 {
		t.Log(stacks[0])
		t.Error(len(stacks[0].Calls))
	}
}

func TestTrimSourcePrefix1(t *testing.T) {
	data, err := os.Open("testdata/data1.txt")
	if err != nil {
		panic(err)
	}
	trace, err := ParseTraceback(data)
	if err != nil {
		panic(err)
	}
	stacks := trace.Stacks
	stacks = ExcludeGotest(stacks)
	stacks = TrimSourcePrefix(stacks)
	if stacks[0].Calls[0].Source != "runtime/sema.goc" {
		t.Fatal(stacks[0].Calls[0])
	}
	if stacks[0].Calls[5].Source != "github.com/t-yuki/mygosandbox/go2qfix/go2qfix_test.go" {
		t.Fatal(stacks[0].Calls[5])
	}
}
