package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/Frank-Mayer/bm/internal/benchmark"
	"github.com/Frank-Mayer/bm/internal/cli"
)

func main() {
	opt, err := cli.Parse(os.Args)
	if err != nil {
		panic(err)
	}

	var cmd *exec.Cmd
	switch len(opt.Command) {
	case 0:
		panic("no command found. Usage: bm [options...] -- [command]")
	case 1:
		cmd = exec.Command(opt.Command[0])
	default:
		cmd = exec.Command(opt.Command[0], opt.Command[1:]...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	addr := benchmark.StartServer()
	url := "http://" + addr
	fmt.Println("Graphical output available at", url)
	if err := open(url); err != nil {
		panic(err)
	}

	if err := benchmark.Run(cmd); err != nil {
		panic(err)
	}

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
