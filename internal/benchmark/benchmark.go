package benchmark

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type Data struct {
	// Title of the measurement
	Title string `json:"title"`
	// Time of measurement
	Time time.Time `json:"time"`
	// Memory percentages
	CPU float32 `json:"cpu"`
	// Memory usage
	Memory float32 `json:"memory"`
}

func Run(cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		return errors.Join(fmt.Errorf("failed to start command: %v", cmd.Args), err)
	}

	title := strings.Join(cmd.Args, " ")

	p, err := process.NewProcess(int32(cmd.Process.Pid))
	if err != nil {
		return errors.Join(fmt.Errorf("failed to get process: %v", cmd.Process.Pid), err)
	}
	for {
		<-time.After(500 * time.Millisecond)
		cpu, err := p.CPUPercent()
		if err != nil {
			break
		}
		mem, err := p.MemoryPercent()
		if err != nil {
			break
		}
		Append(Data{
			Title:  title,
			Time:   time.Now(),
			CPU:    float32(cpu),
			Memory: mem,
		})
	}

	return nil
}

var dataChan = make(chan Data)

func Append(data Data) {
	dataChan <- data
}
