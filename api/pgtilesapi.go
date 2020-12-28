package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ginfolderstructure/lib"
	"github.com/ginfolderstructure/model"
)

func GetTiles(c *gin.Context) {
	service := c.Param("service")
	uuid := c.Param("uuid")
	ts := c.Param("timestamp")
	authHeader := c.GetHeader("Authorization")
	tokenExists, tokenHeader := lib.CheckForBearerToken(authHeader)
	if tokenExists == true {
		c.Set("auth_token", tokenHeader)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized Access",
		})
		c.Abort()
	}
	var (
		z    uint64
		x    uint64
		y    uint64
		err  error
		data []byte
	)
	z, err = strconv.ParseUint(c.Param("z"), 10, 32)
	x, err = strconv.ParseUint(c.Param("x"), 10, 32)
	y, err = strconv.ParseUint(c.Param("y"), 10, 32)

	data, err = model.GetTiles(service, uuid, ts, uint64(z), uint64(x), uint64(y))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"messages": err.Error()})
		return
	}
	if len(data) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"messages": "Tile data is not available."})
		return
	}
	c.Data(http.StatusOK, "application/x-protobuf", data)
}

func GetSummary(c *gin.Context) {
	service := c.Param("service")
	uuid := c.Param("uuid")
	timestamp := c.Param("timestamp")

	var (
		err  error
		data []byte
	)
	data, err = model.GetSummary("", service, uuid, timestamp)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "min max file has not generated for the service"})
	}
	var result []map[string]map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	for key, res := range result {
		fmt.Println("Reading Value for Key :", key)
		if res[timestamp] != nil {
			c.JSON(http.StatusOK, gin.H{"min": res[timestamp]["min"], "max": res[timestamp]["max"]})
			break
		}
	}
}
