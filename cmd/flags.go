package cmd

import (
	"flag"
	"os"
	"strings"
)

var (
	filter *string
	format *string
	help   *bool

	fNrepeat *int
	fProcs   *int
	fNP      *int
)

func init() {
	RegisterFlags(flag.NewFlagSet(os.Args[0], flag.ExitOnError))
}

func RegisterFlags(flags *flag.FlagSet) {
	help = flags.Bool("h", false, "show this help")
	filter = flags.String("f", "trimstd,notest", `stack trace filters by comma-separated list
	trimstd:  exclude GOROOT function calls but leave one
	nostd:    exclude GOROOT function calls completely
	notest:   exclude testing function calls
	top:      remove lower function calls`)
	format = flags.String("t", "col", `output format
	raw: as-is and no filtering
	col: column formatted text
	qfix: vim quickfix output format with errorformat: '%f:%l:\ %m'. you should use with 'nostd,notest,top' filters
	json: JSON format`)

	fNrepeat = flags.Int("n", 1, "repeat the test N times while it passes")
	fProcs = flags.Int("p", 0, "set GOMAXPROCS")
	fNP = flags.Int("np", 0, "similar to a combination of `-n` and `-p` but increment GOMAXPROCS from 1 for each repeat")
}

func ParseFlags(flags *flag.FlagSet) (otherArgs []string) {
	args := make([]string, 0, 10)
	otherArgs = make([]string, 0, 10)

	var skip bool
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if skip {
			otherArgs = append(otherArgs, arg)
			continue
		}
		pair := strings.SplitN(arg, "=", 2)
		key := strings.TrimLeft(pair[0], "-")
		name := strings.TrimRight(key, "0123456789")
		value := strings.TrimPrefix(key, name)
		if value != "" && len(pair) == 1 {
			arg = "-" + name + "=" + value
			pair = []string{name, value}
		}
		switch {
		case flags.Lookup(name) != nil || name == "h" || name == "help":
			args = append(args, arg)
			if i+1 < len(os.Args) && len(pair) == 1 {
				fv, ok := flags.Lookup(name).Value.(interface {
					IsBoolFlag() bool
				})
				if !ok || !fv.IsBoolFlag() {
					args = append(args, os.Args[i+1])
					i++
				}
			}
		case name == "":
			skip = true
		default:
			otherArgs = append(otherArgs, arg)
		}
	}
	flags.Parse(args)
	if *help {
		flags.Usage()
		os.Exit(2)
	}
	return
}
