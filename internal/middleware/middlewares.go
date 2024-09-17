package middleware

import "github.com/gin-gonic/gin"

func AliasTopTours(c *gin.Context) {
	c.Request.URL.Query().Add("limit", "5")
	c.Request.URL.Query().Add("sort", "-ratingsAverage,price")
	c.Request.URL.Query().Add("fields", "name,price,ratingsAverage,summary,difficulty")
	c.Next()
}
