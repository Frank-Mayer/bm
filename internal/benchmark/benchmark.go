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
	CPU float64 `json:"cpu"`
	// Virtual memory size
	VMS uint64 `json:"vms"`
	// Resident memory size
	RSS uint64 `json:"rss"`
	// Stack size
	Stack uint64 `json:"stack"`
	// High Water Mark
	HWM uint64 `json:"hwm"`
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
		mem, err := p.MemoryInfo()
		if err != nil {
			break
		}
		Append(Data{
			Title: title,
			Time:  time.Now(),
			CPU:   cpu,
			VMS:   mem.VMS,
			RSS:   mem.RSS,
			Stack: mem.Stack,
			HWM:   mem.HWM,
		})
	}

	return nil
}

var dataChan = make(chan Data)

func Append(data Data) {
	dataChan <- data
}
