package main

import (
	"log"
	"net/http"
	"os"

	cpuinfo "checker/library/cpu"
	diskinfo "checker/library/disk"
	gpuinfo "checker/library/gpu"
	hostinfo "checker/library/host"
	memoryinfo "checker/library/memory"
	networkinfo "checker/library/network"
	processinfo "checker/library/process"
	sensorinfo "checker/library/sensor"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func configureCors() cors.Config {
	return cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
}

func wsCPUInfoHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade for CPU info: %v", err)
		return
	}
	defer conn.Close()

	for {
		data := cpuinfo.GetCPUInfo
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("Failed to send CPU info over websocket: %v", err)
			break
		}
	}
}

func wsMemoryInfoHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade for Memory info: %v", err)
		return
	}
	defer conn.Close()

	for {
		data := memoryinfo.GetMemoryInfo
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("Failed to send memory info over websocket: %v", err)
			break
		}
	}
}

func wsDiskInfoHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade for Disk info: %v", err)
		return
	}
	defer conn.Close()

	for {
		data := diskinfo.GetDiskInfo
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("Failed to send memory info over websocket: %v", err)
			break
		}
	}
}

func wsNetworkInfoHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade for Network info: %v", err)
		return
	}
	defer conn.Close()

	for {
		data := networkinfo.GetNetworkInfo
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("Failed to send memory info over websocket: %v", err)
			break
		}
	}
}

func wsNetworkPidsInfo(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade for Network Interfaces info: %v", err)
		return
	}
	defer conn.Close()

	for {
		data := networkinfo.GetPids
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("Failed to send network interfaces info over websocket: %v", err)
			break
		}
	}
}

func wsNetworkInterfacesInfo(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade for Network Interfaces info: %v", err)
		return
	}
	defer conn.Close()

	for {
		data := networkinfo.GetInterfaces
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("Failed to send network interfaces info over websocket: %v", err)
			break
		}
	}
}

func wsNetworkConnectionsInfo(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade for Network Connections info: %v", err)
		return
	}
	defer conn.Close()

	for {
		data := networkinfo.GetConnections
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("Failed to send network connections info over websocket: %v", err)
			break
		}
	}
}

func wsNetworkIOCountersInfo(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade for Network IO Counters info: %v", err)
		return
	}
	defer conn.Close()

	for {
		ioCounters, err := networkinfo.GetIOCounters()
		if err != nil {
			log.Printf("Error getting IO counters: %v", err)
			break
		}

		if err := conn.WriteJSON(gin.H{"io_counters": ioCounters}); err != nil {
			log.Printf("Failed to send network IO counters info over websocket: %v", err)
			break
		}
	}
}

func wsNetworkConntrackStatsInfo(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade for Network Conntrack Stats info: %v", err)
		return
	}
	defer conn.Close()

	for {
		conntrackstats, err := networkinfo.GetConntrackStats()
		if err != nil {
			log.Printf("Error getting Conntrack Stats: %v", err)
			break
		}

		if err := conn.WriteJSON(gin.H{"conntrack_stats": conntrackstats}); err != nil {
			log.Printf("Failed to send network conntrack stats info over websocket: %v", err)
			break
		}
	}
}

func wsProcessInfoHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade for Process info: %v", err)
		return
	}
	defer conn.Close()

	for {
		data := processinfo.GetProcessInfo
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("Failed to send memory info over websocket: %v", err)
			break
		}
	}
}

func wsSensorInfoHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade for Sensor info: %v", err)
		return
	}
	defer conn.Close()

	for {
		data := sensorinfo.GetSensorInfo
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("Failed to send memory info over websocket: %v", err)
			break
		}
	}
}

func wsGpuInfoHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade for GPU info: %v", err)
		return
	}
	defer conn.Close()

	for {
		data := gpuinfo.GetGpuInfo
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("Failed to send memory info over websocket: %v", err)
			break
		}
	}
}

func initializeRoutes(r *gin.Engine) {
	metrics := r.Group("/metrics")
	{
		metrics.GET("/system", hostinfo.GetSystemInfo)
		metrics.GET("/cpu", cpuinfo.GetCPUInfo)
		metrics.GET("/memory", memoryinfo.GetMemoryInfo)
		metrics.GET("/disk", diskinfo.GetDiskInfo)
		metrics.GET("/network", networkinfo.GetNetworkInfo)
		metrics.GET("/process", processinfo.GetProcessInfo)
		metrics.GET("/sensors", sensorinfo.GetSensorInfo)
		metrics.GET("/gpu", gpuinfo.GetGpuInfo)
	}

	r.GET("/ws/cpu", wsCPUInfoHandler)
	r.GET("/ws/memory", wsMemoryInfoHandler)
	r.GET("/ws/disk", wsDiskInfoHandler)
	r.GET("/ws/network", wsNetworkInfoHandler)
	r.GET("/ws/network/pids", wsNetworkPidsInfo)
	r.GET("/ws/network/interfaces", wsNetworkInterfacesInfo)
	r.GET("/ws/network/connections", wsNetworkConnectionsInfo)
	r.GET("/ws/network/iocounters", wsNetworkIOCountersInfo)
	r.GET("/ws/network/conntracksstats", wsNetworkConntrackStatsInfo)
	r.GET("/ws/process", wsProcessInfoHandler)
	r.GET("/ws/sensors", wsSensorInfoHandler)
	r.GET("/ws/gpu", wsGpuInfoHandler)
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(cors.New(configureCors()))

	initializeRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "33551"
	}

	log.Printf("Starting server on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
