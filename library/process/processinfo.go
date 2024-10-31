package processinfo

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/process"
)

// GetProcessInfo retrieves process information and returns it as JSON
func GetProcessInfo(c *gin.Context) {
	processInfoList := []map[string]interface{}{}
	processes, err := process.Processes() // Mengambil semua proses yang berjalan
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve processes"})
		return
	}

	// Menggunakan WaitGroup untuk paralelisme
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, proc := range processes {
		wg.Add(1)
		go func(proc *process.Process) {
			defer wg.Done()
			processInfo := getProcessDetails(proc)
			mu.Lock()
			processInfoList = append(processInfoList, processInfo)
			mu.Unlock()
		}(proc)
	}

	wg.Wait()
	c.JSON(http.StatusOK, processInfoList)
}

// getProcessDetails collects detailed information of a single process
func getProcessDetails(proc *process.Process) map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	info := make(map[string]interface{})

	// Mendapatkan berbagai informasi proses
	if pid, err := proc.PpidWithContext(ctx); err == nil {
		info["pid"] = pid
	}
	if name, err := proc.NameWithContext(ctx); err == nil {
		info["name"] = name
	}
	if exe, err := proc.ExeWithContext(ctx); err == nil {
		info["exe"] = exe
	}
	if cmdline, err := proc.CmdlineWithContext(ctx); err == nil {
		info["cmdline"] = cmdline
	}
	if memInfo, err := proc.MemoryInfoWithContext(ctx); err == nil {
		info["memory_info"] = memInfo
	}
	if cpuPercent, err := proc.CPUPercentWithContext(ctx); err == nil {
		info["cpu_percent"] = cpuPercent
	}
	if createTime, err := proc.CreateTimeWithContext(ctx); err == nil {
		info["create_time"] = createTime
	}
	if numThreads, err := proc.NumThreadsWithContext(ctx); err == nil {
		info["num_threads"] = numThreads
	}
	if status, err := proc.StatusWithContext(ctx); err == nil {
		info["status"] = status
	}
	if nice, err := proc.NiceWithContext(ctx); err == nil {
		info["nice"] = nice
	}
	if threads, err := proc.ThreadsWithContext(ctx); err == nil {
		info["threads"] = threads
	}
	if openFiles, err := proc.OpenFilesWithContext(ctx); err == nil {
		info["open_files"] = openFiles
	}
	if children, err := proc.ChildrenWithContext(ctx); err == nil {
		info["children"] = children
	}
	if connections, err := proc.ConnectionsWithContext(ctx); err == nil {
		info["connections"] = connections
	}

	return info
}
