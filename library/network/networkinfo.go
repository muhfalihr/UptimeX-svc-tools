package networkinfo

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/net"
)

func GetPids() ([]int32, error) {
	pids, err := net.Pids()
	if err != nil {
		log.Printf("Error getting PIDs: %v\n", err)
		return nil, err
	}
	return pids, nil
}

func GetInterfaces() ([]net.InterfaceStat, error) {
	addrs, err := net.Interfaces()
	if err != nil {
		log.Printf("Error getting interfaces: %v\n", err)
		return nil, err
	}
	return addrs, nil
}

func GetConnections() ([]net.ConnectionStat, error) {
	connections, err := net.Connections("inet")
	if err != nil {
		log.Printf("Error getting connections: %v\n", err)
		return nil, err
	}
	return connections, nil
}

func GetIOCounters() ([]net.IOCountersStat, error) {
	ioCounters, err := net.IOCounters(true)
	if err != nil {
		log.Printf("Error getting IO counters: %v\n", err)
		return nil, err
	}
	return ioCounters, nil
}

func GetConntrackStats() ([]net.ConntrackStat, error) {
	conntrackStats, err := net.ConntrackStats(true)
	if err != nil {
		log.Printf("Error getting conntrack stats: %v\n", err)
		return nil, err
	}
	return conntrackStats, nil
}

func GetNetworkInfo(c *gin.Context) {
	var wg sync.WaitGroup
	infoCh := make(chan map[string]interface{})
	networkInfo := make([]map[string]interface{}, 0)

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

	go func() {
		wg.Wait()
		close(infoCh)
	}()

	for info := range infoCh {
		networkInfo = append(networkInfo, info)
	}

	c.JSON(http.StatusOK, networkInfo)
}
