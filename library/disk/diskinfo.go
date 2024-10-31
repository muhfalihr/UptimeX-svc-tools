package diskinfo

import (
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/disk"
)

func GetDiskInfo(c *gin.Context) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // Set a timeout of 5 seconds
	defer cancel()

	var wg sync.WaitGroup
	diskInfo := make([]gin.H, 0, len(partitions))
	infoCh := make(chan gin.H)

	for _, partition := range partitions {
		wg.Add(1)

		go func(partition disk.PartitionStat) {
			defer wg.Done()

			usage, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				if os.IsPermission(err) {
					return
				}
				return
			}

			ioCounters, _ := disk.IOCountersWithContext(ctx, partition.Device)
			ioCounter := ioCounters[partition.Device]

			label := partition.Fstype
			serialNumber := partition.Device

			infoCh <- gin.H{
				"device":         partition.Device,
				"mountpoint":     partition.Mountpoint,
				"filesystem":     partition.Fstype,
				"total_space":    usage.Total,
				"used_space":     usage.Used,
				"free_space":     usage.Free,
				"used_percent":   usage.UsedPercent,
				"io_read_count":  ioCounter.ReadCount,
				"io_write_count": ioCounter.WriteCount,
				"io_read_bytes":  ioCounter.ReadBytes,
				"io_write_bytes": ioCounter.WriteBytes,
				"label":          label,
				"serial_number":  serialNumber,
			}
		}(partition)
	}

	go func() {
		wg.Wait()
		close(infoCh)
	}()

	for info := range infoCh {
		diskInfo = append(diskInfo, info)
	}

	c.JSON(http.StatusOK, diskInfo)
}
