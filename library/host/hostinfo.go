package systeminfo

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/host"
)

func GetSystemInfo(c *gin.Context) {
	sysInfoChan := make(chan *host.InfoStat)
	errChan := make(chan error)

	go func() {
		sysInfo, err := host.Info()
		if err != nil {
			errChan <- err
			return
		}
		sysInfoChan <- sysInfo
	}()

	select {
	case sysInfo := <-sysInfoChan:
		bootTime, _ := host.BootTime()
		uptime, _ := host.Uptime()
		users, _ := host.Users()
		kernelArch, _ := host.KernelArch()
		kernelVersion, _ := host.KernelVersion()
		hostID, _ := host.HostID()

		c.JSON(http.StatusOK, gin.H{
			"system":         sysInfo.OS,
			"hostname":       sysInfo.Hostname,
			"platform":       sysInfo.Platform,
			"version":        sysInfo.PlatformVersion,
			"arch":           runtime.GOARCH,
			"boot_time":      bootTime,
			"uptime":         uptime,
			"users":          users,
			"kernel_arch":    kernelArch,
			"kernel_version": kernelVersion,
			"host_id":        hostID,
		})
	case err := <-errChan:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
