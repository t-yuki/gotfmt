package cmd_test

import (
	"flag"
	"testing"

	"github.com/t-yuki/gotfmt/cmd"
)

func TestParseFlags(t *testing.T) {
	var flags = flag.NewFlagSet("test", flag.ExitOnError)
	cmd.RegisterFlags(flags)
	others := cmd.ParseFlags(flags, []string{"-t=json", "otherarg"})
	if val := flags.Lookup("t").Value.String(); val != "json" {
		t.Fatal("want: json but:", val)
	}
	if others[0] != "otherarg" {
		t.Fatal(others[0])
	}
}

func TestParseFlags_2(t *testing.T) {
	var flags = flag.NewFlagSet("test", flag.ExitOnError)
	cmd.RegisterFlags(flags)
	others := cmd.ParseFlags(flags, []string{"-t", "json", "otherarg"})
	if val := flags.Lookup("t").Value.String(); val != "json" {
		t.Fatal("want: json but:", val)
	}
	if others[0] != "otherarg" {
		t.Fatal(others[0])
	}
}

func TestParseFlags_3(t *testing.T) {
	var flags = flag.NewFlagSet("test", flag.ExitOnError)
	cmd.RegisterFlags(flags)
	others := cmd.ParseFlags(flags, []string{"-np3", "otherarg"})
	if val := flags.Lookup("np").Value.String(); val != "3" {
		t.Fatal(val)
	}
	if others[0] != "otherarg" {
		t.Fatal(others[0])
	}
}

func TestParseFlags_skip(t *testing.T) {
	var flags = flag.NewFlagSet("test", flag.ExitOnError)
	cmd.RegisterFlags(flags)
	others := cmd.ParseFlags(flags, []string{"-", "-t"})
	if others[0] != "-t" {
		t.Fatal(others[0])
	}
}
