package cmd_test

import (
	"flag"
	"os"
	"testing"

	"github.com/t-yuki/gotfmt/cmd"
)

func testMain(t *testing.T, args ...string) {
	if os.Getenv(".gotfmt.TestMain") == "" {
		os.Setenv(".gotfmt.TestMain", "1")
		defer os.Setenv(".gotfmt.TestMain", "")
		runMain(t, append(args, "-np0")...)
		runMain(t, args...)
		runMain(t, append(args, "-run=TestMain")...)
		runMain(t, append(args, "-test.run=TestMain")...)
		runMain(t, append(args, "-np1")...)
		runMain(t, append(args, "-np2")...)
		runMain(t, append(args, "-np1", "-run=TestMain")...)
		runMain(t, append(args, "-np2", "-run=TestMain")...)
		runMain(t, append(args, "-np1", "-test.run=TestMain")...)
		runMain(t, append(args, "-np2", "-test.run=TestMain")...)
	}
}

func runMain(t *testing.T, args ...string) {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cmd.RegisterFlags(flags)
	args = cmd.ParseFlags(flags, args)
	cmd.Main(args)
}

func TestMain_test(t *testing.T) {
	testMain(t, "test")
}

func TestMain_race(t *testing.T) {
	testMain(t, "test", "-race")
}

func TestMain_cover(t *testing.T) {
	testMain(t, "test", "-cover")
}
