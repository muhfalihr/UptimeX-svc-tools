package memoryinfo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/mem"
)

func GetMemoryInfo(c *gin.Context) {
	memInfoChan := make(chan *mem.VirtualMemoryStat)
	swapInfoChan := make(chan *mem.SwapMemoryStat)
	errChan := make(chan error)

	go func() {
		memInfo, err := mem.VirtualMemory()
		if err != nil {
			errChan <- err
			return
		}
		memInfoChan <- memInfo
	}()

	go func() {
		swapInfo, err := mem.SwapMemory()
		if err != nil {
			errChan <- err
			return
		}
		swapInfoChan <- swapInfo
	}()

	select {
	case memInfo := <-memInfoChan:
		select {
		case swapInfo := <-swapInfoChan:
			c.JSON(http.StatusOK, gin.H{
				"total_memory":        memInfo.Total,
				"available_memory":    memInfo.Available,
				"used_memory":         memInfo.Used,
				"free_memory":         memInfo.Free,
				"used_memory_percent": memInfo.UsedPercent,
				"total_swap":          swapInfo.Total,
				"used_swap":           swapInfo.Used,
				"free_swap":           swapInfo.Free,
				"used_swap_percent":   swapInfo.UsedPercent,
			})
		case err := <-errChan:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	case err := <-errChan:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
