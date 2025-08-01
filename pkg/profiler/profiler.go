package profiler

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

type Profiler struct {
	fileCPU *os.File
	fileMem *os.File
	isOn    bool
}

func New(isOn bool, cpuProfile, memProfile string) (*Profiler, error) {
	if !isOn {
		return &Profiler{
			isOn: isOn,
		}, nil
	}
	if err := os.MkdirAll("profiles", 0o755); err != nil {
		return nil, err
	}
	// создаём файл журнала профилирования cpu
	fcpu, err := os.Create(`profiles/` + cpuProfile)
	if err != nil {
		return nil, err
	}
	fmem, err := os.Create(`profiles/` + memProfile)
	if err != nil {
		return nil, err
	}
	return &Profiler{
		fileCPU: fcpu,
		fileMem: fmem,
		isOn:    isOn,
	}, nil
}

func (p *Profiler) Close() error {
	if !p.isOn {
		return nil
	}
	if err := p.fileCPU.Close(); err != nil {
		return fmt.Errorf("error in Profiler Close - fileCPU.Close: %v", err)
	}

	if err := p.fileMem.Close(); err != nil {
		return fmt.Errorf("error in Profiler Close - fileMem.Close: %v", err)
	}
	return nil
}

func (p *Profiler) Start() {
	if !p.isOn {
		return
	}
	if err := pprof.StartCPUProfile(p.fileCPU); err != nil {
		panic(err)
	}
	time.AfterFunc(30*time.Second, p.stop)
}

func (p *Profiler) stop() {
	pprof.StopCPUProfile()
	runtime.GC()
	if err := pprof.WriteHeapProfile(p.fileMem); err != nil {
		panic(err)
	}
}
