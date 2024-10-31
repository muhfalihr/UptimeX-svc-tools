package cpuinfo

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
)

func GetCPUInfo(c *gin.Context) {
	var wg sync.WaitGroup
	cpuInfoChan := make(chan []cpu.InfoStat)
	cpuPercentChan := make(chan []float64)
	cpuTimesChan := make(chan []cpu.TimesStat)
	errChan := make(chan error, 3)

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		info, err := cpu.Info()
		if err != nil {
			errChan <- err
			return
		}
		cpuInfoChan <- info
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		percent, err := cpu.Percent(0, true)
		if err != nil {
			errChan <- err
			return
		}
		cpuPercentChan <- percent
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		times, err := cpu.TimesWithContext(ctx, true)
		if err != nil {
			errChan <- err
			return
		}
		cpuTimesChan <- times
	}()

	go func() {
		wg.Wait()
		close(cpuInfoChan)
		close(cpuPercentChan)
		close(cpuTimesChan)
		close(errChan)
	}()

	var err error
	select {
	case err = <-errChan:
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case cpuInfo := <-cpuInfoChan:
		var cpuPercent []float64
		select {
		case cpuPercent = <-cpuPercentChan:
		case err = <-errChan:
		}
		var cpuTimes []cpu.TimesStat
		select {
		case cpuTimes = <-cpuTimesChan:
		case err = <-errChan:
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logicalCPUCount, err := cpu.Counts(true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		physicalCPUCount, err := cpu.Counts(false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"cpu_info":             cpuInfo,
			"cpu_count_physical":   physicalCPUCount,
			"cpu_count_logical":    logicalCPUCount,
			"cpu_percent_per_core": cpuPercent,
			"cpu_times":            cpuTimes,
		})
	}
}
