package cli

import (
	"fmt"
)

type Options struct {
	Command []string
}

func Parse(args []string) (*Options, error) {
	if len(args) > 1 {
		return &Options{Command: args[1:]}, nil
	}
	return nil, fmt.Errorf("no command found. Usage: %s [command]", args[0])
}
