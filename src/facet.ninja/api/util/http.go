package util

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const GET = "GET"
const POST = "POST"
const NOT_FOUND = "NOT_FOUND"

func SetResponseCode(result interface{}, error error, context *gin.Context) {
	SetCorsHeaders(context)
	if error != nil {
		log.Print(error)
		if error.Error() == NOT_FOUND {
			context.JSON(http.StatusNotFound, error)
		} else {
			context.JSON(http.StatusInternalServerError, error.Error())
		}

	} else if context.Request.Method == POST {
		context.JSON(http.StatusCreated, result)
	} else {
		context.JSON(http.StatusOK, result)
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
