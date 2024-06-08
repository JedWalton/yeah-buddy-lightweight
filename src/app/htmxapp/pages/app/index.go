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
	const N = 20
	arcsData := make([]map[string]interface{}, N)
	for i := 0; i < N; i++ {
		arcsData[i] = map[string]interface{}{
			"startLat": (randomSource.Float64() - 0.5) * 180,
			"startLng": (randomSource.Float64() - 0.5) * 360,
			"endLat":   (randomSource.Float64() - 0.5) * 180,
			"endLng":   (randomSource.Float64() - 0.5) * 360,
			"color": []string{
				[]string{"red", "white", "blue", "green"}[randomSource.Intn(4)],
				[]string{"red", "white", "blue", "green"}[randomSource.Intn(4)],
			},
		}
	}
	c.JSON(http.StatusOK, arcsData)
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
