package util

import (
	"github.com/gin-gonic/gin"
	"log"
)

const GET = "GET"
const POST = "POST"
const DELETE = "DELETE"
const OPTIONS = "OPTIONS"
const NOT_FOUND = "NOT_FOUND"

func SetResponseCode(result interface{}, error error, c *gin.Context) {
	SetCorsHeaders(c)
	if error != nil {
		log.Print(error)
		if error.Error() == NOT_FOUND {
			c.JSON(404, error)
		} else {
			c.JSON(500, error.Error())
		}

	} else if c.Request.Method == POST {
		c.JSON(201, result)
	} else {
		c.JSON(200, result)
	}
}

func SetCorsHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
}

func Options(c *gin.Context) {
	SetResponseCode(nil, nil, c)
}
