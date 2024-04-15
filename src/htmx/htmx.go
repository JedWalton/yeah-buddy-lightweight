package htmx

import (
	"github.com/gin-gonic/gin"
	"sync"
)

var (
	items = []string{} // Slice to store items
	mutex sync.Mutex   // Mutex to handle concurrent modifications
)

func Init(r *gin.Engine) {
	// Serve HTML file
	r.GET("/", func(c *gin.Context) {
		c.File("./htmx/templates/index.html")
	})

	// HTMX endpoint
	r.GET("/data", func(c *gin.Context) {
		c.String(200, "Hello from HTMX and Gin!")
	})

	r.GET("/items", func(c *gin.Context) {
		c.HTML(200, "items.html", gin.H{
			"items": items,
		})
	})

	r.POST("/add-item", func(c *gin.Context) {
		item := c.PostForm("item")
		mutex.Lock()
		items = append(items, item)
		mutex.Unlock()
		c.Redirect(303, "/items")
	})

	r.POST("/delete-item", func(c *gin.Context) {
		index := c.PostForm("index")
		i := convertToInt(index)
		mutex.Lock()
		items = append(items[:i], items[i+1:]...)
		mutex.Unlock()
		c.Redirect(303, "/items")
	})

	r.LoadHTMLGlob("./htmx/templates/*")
}

func convertToInt(s string) int {
	// conversion logic here
	return 0
}
