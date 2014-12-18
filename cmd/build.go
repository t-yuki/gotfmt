package cmd

import (
	"errors"
	"os/exec"
)

func build(output string, args []string) error {
	args = extractBuildFlags(args)
	args = append([]string{"test", "-c", "-o", output}, args...)
	cmd := exec.Command("go", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) != 0 {
			return errors.New(string(out))
		}
		return err
	}
	return nil
}

func extractBuildFlags(args []string) []string {
	ret := make([]string, 0, len(args))
	for _, arg := range args {
		switch arg {
		case "-race", "-cover", "-covermode", "-coverpkg":
			break
		default:
			continue
		}
		ret = append(ret, arg)
	}
	return ret

}
