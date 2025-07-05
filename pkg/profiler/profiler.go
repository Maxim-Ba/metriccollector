package profiler

import (
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

type Profiler struct {
	fileCPU *os.File
	fileMem *os.File
}

func New(cpuProfile, memProfile string) *Profiler {
	if err := os.MkdirAll("profiles", 0o755); err != nil {
		panic(err)
	}
	// создаём файл журнала профилирования cpu
	fcpu, err := os.Create(`profiles/` + cpuProfile)
	if err != nil {
		panic(err)
	}
	fmem, err := os.Create(`profiles/` +memProfile)
	if err != nil {
		panic(err)
	}
	return &Profiler{
		fileCPU: fcpu,
		fileMem: fmem,
	}
}

func (p *Profiler) Close() {
	p.fileCPU.Close()
	p.fileMem.Close()
}

func (p *Profiler) Start() {
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
