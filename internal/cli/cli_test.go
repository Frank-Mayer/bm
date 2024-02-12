package cli_test

import (
	"fmt"
	"testing"

	"github.com/Frank-Mayer/bm/internal/cli"
)

func TestCli(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in  []string
		out cli.Options
	}{
		{[]string{"bm", "echo", "hello"}, cli.Options{Command: []string{"echo", "hello"}}},
		{[]string{"bm", "sleep", "10"}, cli.Options{Command: []string{"sleep", "10"}}},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			t.Parallel()
			o, err := cli.Parse(c.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(o.Command) != len(c.out.Command) {
				t.Errorf("expected command %v, got %v", c.out.Command, o.Command)
			} else {
				for i, v := range o.Command {
					if v != c.out.Command[i] {
						t.Errorf("expected command %v, got %v", c.out.Command, o.Command)
					}
				}
			}
		})
	}
}
