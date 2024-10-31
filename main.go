package main

import (
	"log"
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
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "33551"
	}

	log.Printf("Starting server on port %s...", port)
	if err := r.Run(":33551"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
