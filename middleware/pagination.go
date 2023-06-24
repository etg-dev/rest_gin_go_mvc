package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Paginate() gin.HandlerFunc {
	return func(c *gin.Context) {

		pageStr := c.DefaultQuery("page", "1")
		pageSizeStr := c.DefaultQuery("pageSize", "10")

		pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize"})
			return
		}

		page, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
			return
		}

		c.Set("page", page-1)
		c.Set("pageSize", pageSize)

		c.Next()
	}
}
