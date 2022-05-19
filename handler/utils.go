package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func badRequestWithError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status": "fail",
		"error":  err.Error(),
	})
}

func serverError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status": "fail",
		"error":  err.Error()})
}

func statusOkWithSuccess(c *gin.Context, mapValues map[string]interface{}, structValues interface{}) {
	h := gin.H{"status": "success"}

	if mapValues == nil && structValues == nil {
		c.JSON(http.StatusOK, h)
	} else if structValues == nil {
		for k, v := range mapValues {
			h[k] = v
		}
		c.JSON(http.StatusOK, h)
	} else {
		if err := json.NewEncoder(c.Writer).Encode(structValues); err != nil {
			serverError(c, err)
		}
	}

}
