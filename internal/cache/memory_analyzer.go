package cache

import (
	"runtime"

	"sync"

	"code.cloudfoundry.org/log-cache/internal/metrics"
	sigar "github.com/cloudfoundry/gosigar"
)

// MemoryAnalyzer reports the available and total memory.
type MemoryAnalyzer struct {
	// metrics
	setAvail func(value float64)
	setTotal func(value float64)
	setHeap  func(value float64)

	// cache
	avail uint64
	total uint64
	heap  uint64

	sync.Mutex
}

// NewMemoryAnalyzer creates and returns a new MemoryAnalyzer.
func NewMemoryAnalyzer(metrics metrics.Initializer) *MemoryAnalyzer {
	return &MemoryAnalyzer{
		setAvail: metrics.NewGauge("log_cache_available_system_memory", "bytes"),
		setHeap:  metrics.NewGauge("log_cache_heap_in_use_memory", "bytes"),
		setTotal: metrics.NewGauge("log_cache_total_system_memory", "bytes"),
	}
}

// Memory returns the heap memory and total system memory.
func (a *MemoryAnalyzer) Memory() (heapInUse, available, total uint64) {
	a.Lock()
	defer a.Unlock()

	var m sigar.Mem
	m.Get()

	a.avail = m.ActualFree
	a.total = m.Total

	a.setAvail(float64(m.ActualFree))
	a.setTotal(float64(m.Total))

	var rm runtime.MemStats
	runtime.ReadMemStats(&rm)

	a.heap = rm.HeapInuse
	a.setHeap(float64(a.heap))

	return a.heap, a.avail, a.total
}
