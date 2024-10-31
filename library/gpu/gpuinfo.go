package gpuinfo

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

type NvResult struct {
	Free  uint64 `json:"free"`
	Used  uint64 `json:"used"`
	Total uint64 `json:"total"`
}

type GpuInfo struct {
	NvidiaMemory NvResult `json:"nvidia_memory"`
}

func parseNvidiaSmiOutput(output string) (NvResult, error) {
	lines := strings.Split(output, "\n")
	var free, used, total uint64

	for _, line := range lines {
		if strings.Contains(line, "MiB") && strings.Contains(line, "/") {
			fields := strings.Fields(line)
			totalStr := strings.TrimSuffix(fields[8], "MiB")
			usedStr := strings.TrimSuffix(fields[5], "MiB")
			freeStr := strings.TrimSuffix(fields[10], "MiB")

			var err error
			if total, err = strconv.ParseUint(totalStr, 10, 64); err != nil {
				return NvResult{}, err
			}
			if used, err = strconv.ParseUint(usedStr, 10, 64); err != nil {
				return NvResult{}, err
			}
			if free, err = strconv.ParseUint(freeStr, 10, 64); err != nil {
				return NvResult{}, err
			}
			break
		}
	}
	return NvResult{Free: free, Used: used, Total: total}, nil
}

func pollNvidiaGpuMemory() (NvResult, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=memory.free,memory.used,memory.total", "--format=csv,noheader,nounits")
	output, err := cmd.Output()
	if err != nil {
		return NvResult{}, err
	}
	return parseNvidiaSmiOutput(string(output))
}

func channelWriter(ctx context.Context, wg *sync.WaitGroup, ch chan<- GpuInfo) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:

			if nvidiaMemory, err := pollNvidiaGpuMemory(); err == nil {
				gpuInfo := GpuInfo{
					NvidiaMemory: nvidiaMemory,
				}
				ch <- gpuInfo
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func GetGpuInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch := make(chan GpuInfo)
	var wg sync.WaitGroup
	wg.Add(1)

	go channelWriter(ctx, &wg, ch)

	var results []GpuInfo
	for info := range ch {
		results = append(results, info)
		if len(results) >= 5 {
			break
		}
	}

	wg.Wait()

	c.JSON(http.StatusOK, results)
}
