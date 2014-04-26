package traceback

import "testing"

func TestExcludeGotest(t *testing.T) {
	trace := &Traceback{Stacks: []Stack{
		{Calls: []Call{
			{Source: "testing/testing.go"},
			{Source: "_test/_testmain.go"},
		}},
		{Calls: []Call{
			{Source: "testing/testing.go"},
			{Source: "time/sleep.go"},
		}},
		{ID: 1, Calls: []Call{
			{Source: "mypackage/main_test.go"},
			{Source: "testing/testing.go"},
		}},
		{ID: 3, Calls: []Call{
			{Source: "mypackage/main_test.go"},
			{Source: "time/sleep.go"},
		}},
	}}
	stacks := ExcludeGotest(trace.Stacks)
	if len(stacks) != 2 {
		t.Error(stacks)
	}
	if stacks[0].ID != 1 || len(stacks[0].Calls) != 1 {
		t.Error(stacks[0])
	}
	if stacks[1].ID != 3 || len(stacks[1].Calls) != 2 {
		t.Error(stacks[1])
	}
}

func TestExcludeGoroot(t *testing.T) {
	trace := &Traceback{Stacks: []Stack{
		{ID: 1, Calls: []Call{
			{Source: "/usr/local/go/src/pkg/runtime/time.goc"},
			{Source: "main_test.go"},
			{Source: "/usr/local/go/src/pkg/testing/testing.go"},
		}},
		{ID: 3, Calls: []Call{
			{Source: "/usr/local/go/src/pkg/testing/testing.go"},
			{Source: "/usr/local/go/src/pkg/runtime/runtime.go"},
			{Source: "/usr/local/go/src/pkg/runtime/time.goc"},
		}},
	}}

	stacks := ExcludeGoroot(trace.Stacks, false)
	if len(stacks) != 1 {
		t.Error(stacks)
	}
	if stacks[0].ID != 1 ||
		len(stacks[0].Calls) != 2 ||
		stacks[0].Calls[0].Source != "main_test.go" ||
		stacks[0].Calls[1].Source != "/usr/local/go/src/pkg/testing/testing.go" {
		t.Error(stacks[0])
	}

	stacks = ExcludeGoroot(trace.Stacks, true)
	if len(stacks) != 1 {
		t.Error(stacks)
	}
	if stacks[0].ID != 1 ||
		len(stacks[0].Calls) != 3 ||
		stacks[0].Calls[0].Source != "/usr/local/go/src/pkg/runtime/time.goc" ||
		stacks[0].Calls[1].Source != "main_test.go" ||
		stacks[0].Calls[2].Source != "/usr/local/go/src/pkg/testing/testing.go" {
		t.Error(stacks[0])
	}
}

func TestTrimSourcePrefix(t *testing.T) {
	trace := &Traceback{Stacks: []Stack{
		{ID: 1, Calls: []Call{
			{Source: "/usr/local/go/src/pkg/root.goc"},
			{Source: "/mygoroot/src/pkg/root.go"},
			{Source: "/mygopath/src/main.go"},
		}},
	}}
	stacks := trace.Stacks
	stacks = TrimSourcePrefix(stacks)
	if stacks[0].Calls[0].Source != "root.goc" {
		t.Fatal(stacks[0].Calls[0])
	}
	if stacks[0].Calls[1].Source != "root.go" {
		t.Fatal(stacks[0].Calls[1])
	}
	if stacks[0].Calls[2].Source != "main.go" {
		t.Fatal(stacks[0].Calls[1])
	}
}
