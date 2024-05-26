package app

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

func Init(r *gin.Engine) {
	// Serve the main page
	r.GET("/app", func(c *gin.Context) {
		c.HTML(http.StatusOK, "app.html", nil)
	})

	r.GET("/app2", func(c *gin.Context) {
		c.HTML(http.StatusOK, "app2.html", nil)
	})

	r.GET("/app/graph-data", func(c *gin.Context) {
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

/*
* Only write to the DB every 10 mins to reduce load on the db.
* If is all of last 10 mins, write to db for persistent log entry.


* Provide a calculation of average response
* Implement pings from different regions.
* Implement a way to notify the user if the site is down.
* Generate a graph that displays the uptime for last 7 days.
* Average response time from each region.
* Implement a driver that will coordinate pings from all regions.
* Graph like this that plots avg response time from each region to a specific endpoint.
* * https://nextjs-demo.tailadmin.com/dashboard/stocks
 */
