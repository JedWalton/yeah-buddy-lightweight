package htmx

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	items = make(map[string][]string) // Map to store items by categories
	mutex sync.Mutex                  // Mutex to handle concurrent modifications
)

func Init(r *gin.Engine) {
	// Serve HTML file
	r.GET("/", func(c *gin.Context) {
		c.File("./htmx/templates/index.html")
	})

	r.GET("/items", func(c *gin.Context) {
		category := c.Query("category")
		if category == "" {
			category = "General"
		}
		mutex.Lock()
		itemList, ok := items[category]
		mutex.Unlock()
		if !ok {
			itemList = []string{} // Ensure there is always a slice to pass to the template
		}
		c.HTML(200, "items.html", gin.H{
			"category": category,
			"items":    itemList,
		})
	})

	r.POST("/add-item", func(c *gin.Context) {
		category := c.PostForm("category")
		item := c.PostForm("item")
		mutex.Lock()
		items[category] = append(items[category], item)
		mutex.Unlock()
		c.Redirect(303, "/items?category="+category)
	})
	r.POST("/delete-item", func(c *gin.Context) {
		category := c.PostForm("category")
		indexStr := c.PostForm("index")
		index, err := convertToInt(indexStr)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid index: %v", err)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		if index < 0 || index >= len(items[category]) {
			c.String(http.StatusBadRequest, "Index out of range")
			return
		}

		// Perform the slice operation
		items[category] = append(items[category][:index], items[category][index+1:]...)
		c.Redirect(303, "/items?category="+category)
	})

	r.LoadHTMLGlob("./htmx/templates/*")
}

func convertToInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}
