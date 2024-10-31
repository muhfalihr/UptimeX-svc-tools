package networkinfo

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/net"
)

func GetNetworkInfo(c *gin.Context) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a timeout of 5 seconds
	// defer cancel()

	var wg sync.WaitGroup
	infoCh := make(chan map[string]interface{})
	networkInfo := make([]map[string]interface{}, 0)

	// Goroutine untuk mengambil PID
	wg.Add(1)
	go func() {
		defer wg.Done()
		pids, err := net.Pids()
		if err != nil {
			log.Printf("Error getting PIDs: %v\n", err)
			return
		}
		infoCh <- map[string]interface{}{"pids": pids}
	}()

	// Goroutine untuk mengambil alamat interface jaringan
	wg.Add(1)
	go func() {
		defer wg.Done()
		addrs, err := net.Interfaces()
		if err != nil {
			log.Printf("Error getting interfaces: %v\n", err)
			return
		}
		infoCh <- map[string]interface{}{"interfaces": addrs}
	}()

	// Goroutine untuk mengambil koneksi aktif
	wg.Add(1)
	go func() {
		defer wg.Done()
		connections, err := net.Connections("inet")
		if err != nil {
			log.Printf("Error getting connections: %v\n", err)
			return
		}
		infoCh <- map[string]interface{}{"connections": connections}
	}()

	// Goroutine untuk mengambil IO counters
	wg.Add(1)
	go func() {
		defer wg.Done()
		ioCounters, err := net.IOCounters(false)
		if err != nil {
			log.Printf("Error getting IO counters: %v\n", err)
			return
		}
		infoCh <- map[string]interface{}{"io_counters": ioCounters}
	}()

	// Goroutine untuk mengambil conntrack stats
	wg.Add(1)
	go func() {
		defer wg.Done()
		conntrackStats, err := net.ConntrackStats(true)
		if err != nil {
			log.Printf("Error getting conntrack stats: %v\n", err)
			return
		}
		infoCh <- map[string]interface{}{"conntrack_stats": conntrackStats}
	}()

	// Menutup channel setelah semua goroutine selesai
	go func() {
		wg.Wait()
		close(infoCh)
	}()

	// Mengumpulkan hasil dari setiap goroutine
	for info := range infoCh {
		networkInfo = append(networkInfo, info)
	}

	// Mengirim data JSON ke klien
	c.JSON(http.StatusOK, networkInfo)
}
