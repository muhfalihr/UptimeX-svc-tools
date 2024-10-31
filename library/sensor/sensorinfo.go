package sensorinfo

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/host"
)

type SensorData struct {
	SensorTemperatures []host.TemperatureStat `json:"sensor_temperatures"`
	TemperatureStat    []host.TemperatureStat `json:"temperature_stat"`
}

func GetSensorInfo(c *gin.Context) {
	resultChan := make(chan interface{}, 2)
	var wg sync.WaitGroup
	timeout := time.After(3 * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		temps, err := host.SensorsTemperatures()
		if err != nil {
			resultChan <- map[string]string{"error": "failed to get sensor temperatures"}
			return
		}
		resultChan <- temps
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tempStat, err := host.SensorsTemperatures()
		if err != nil {
			resultChan <- map[string]string{"error": "failed to get temperature stat"}
			return
		}
		resultChan <- tempStat
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	sensorData := SensorData{}
	for {
		select {
		case res, ok := <-resultChan:
			if !ok {
				goto RESULT
			}
			switch v := res.(type) {
			case []host.TemperatureStat:
				if len(sensorData.SensorTemperatures) == 0 {
					sensorData.SensorTemperatures = v
				} else {
					sensorData.TemperatureStat = v
				}
			case map[string]string:
				c.JSON(http.StatusInternalServerError, v)
				return
			}
		case <-timeout:
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timed out"})
			return
		}
	}

RESULT:
	c.JSON(http.StatusOK, sensorData)
}
