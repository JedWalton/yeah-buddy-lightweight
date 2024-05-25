package app

import (
	"github.com/gin-gonic/gin"
	"i-couldve-got-six-reps/api/db"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func Init(r *gin.Engine) {
	// Serve the main page
	r.GET("/app", func(c *gin.Context) {
		c.HTML(http.StatusOK, "app.html", nil)
	})

	r.GET("/latest-uptime-log", func(c *gin.Context) {
		getLatestUptimeLog(c)
	})

	r.GET("/graph-data", func(c *gin.Context) {
		getGraphData(c)
	})
}

func getGraphData(c *gin.Context) {
	var randomSource = rand.New(rand.NewSource(time.Now().UnixNano()))
	data := map[string]float64{
		"value": randomSource.Float64() * 100, // Use the custom random source
	}
	c.JSON(http.StatusOK, data)
}

type UptimeLog struct {
	EndpointID   int       `json:"endpoint_id"`
	StatusCode   int       `json:"status_code"`
	ResponseTime int       `json:"response_time"`
	IsUp         bool      `json:"is_up"`
	Timestamp    time.Time `json:"timestamp"`
}

func getLatestUptimeLog(c *gin.Context) {
	db := db.Init()
	defer db.Close()
	var uptimeLog UptimeLog
	query := `SELECT endpoint_id, status_code, response_time, is_up, timestamp 
              FROM UptimeLogs 
              ORDER BY timestamp DESC 
              LIMIT 1`
	err := db.QueryRow(query).Scan(&uptimeLog.EndpointID, &uptimeLog.StatusCode, &uptimeLog.ResponseTime, &uptimeLog.IsUp, &uptimeLog.Timestamp)
	if err != nil {
		log.Printf("Error fetching log: %v", err) // Add this line for logging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch data"})
		return
	}

	log.Printf("Fetched log: %+v", uptimeLog) // Add this line for logging
	c.JSON(http.StatusOK, uptimeLog)
}
