package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SelectFields(fields []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var selectedFields []string
		if selectStr, ok := c.GetQuery("select"); ok {
			selectedFields = strings.Split(selectStr, ",")
		}
		for _, field := range selectedFields {
			if !contains(fields, field) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field selected"})
				c.Abort()
				return
			}
		}
		c.Set("selectedFields", selectedFields)
		c.Next()
	}
}

func contains(arr []string, val string) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}
