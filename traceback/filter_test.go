package traceback

import (
	"os"
	"testing"
)

func TestFilterGotest1(t *testing.T) {
	data, err := os.Open("testdata/data1.txt")
	if err != nil {
		panic(err)
	}
	stacks, err := ParseStacks(data)
	if err != nil {
		panic(err)
	}
	stacks = FilterGotest(stacks)
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

func TestFilterGotest2(t *testing.T) {
	data, err := os.Open("testdata/data2.txt")
	if err != nil {
		panic(err)
	}
	stacks, err := ParseStacks(data)
	if err != nil {
		panic(err)
	}
	stacks = FilterGotest(stacks)
	if len(stacks) != 0 {
		t.Error(stacks)
	}
}
