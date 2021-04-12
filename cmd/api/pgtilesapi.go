package api

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/local-server-test1/cmd/lib"
	"github.com/local-server-test1/cmd/model"
)

// GetTiles api to get the tile data from postgres
func GetTiles(c *gin.Context) {
	service := c.Param("service")
	uuid := c.Param("uuid")
	ts := c.Param("timestamp")

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

	data, err = model.GetTiles(service, uuid, ts, uint64(z), uint64(x), uint64(y), "", "")
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"messages": err.Error()})
	} else if len(data) == 0 || err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"messages": "Tile data is not available."})
	} else {
		c.Data(http.StatusOK, "application/x-protobuf", data)
	}
}

// GetLDMTiles api to get the tile data from postgres
func GetLDMTiles(c *gin.Context) {
	service := c.MustGet("service").(string)
	uuid := c.Param("uuid")
	ts := c.Param("timestamp")
	dbtype := c.Param("type")
	lyrtype := c.Param("lyrtype")

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

	data, err = model.GetTiles(service, uuid, ts, uint64(z), uint64(x), uint64(y), dbtype, lyrtype)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"messages": err.Error()})
	} else if len(data) == 0 || err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"messages": "Tile data is not available."})
	} else {
		c.Data(http.StatusOK, "application/x-protobuf", data)
	}
}

// GetSummary api to get service summary
func GetSummary(c *gin.Context) {
	service := c.Param("service")
	uuid := c.Param("uuid")
	timestamp := c.Param("timestamp")

	var (
		err    error
		minVal float64
		maxVal float64
	)

	minVal, maxVal, err = model.GetSummary(service, uuid, timestamp, "", "")
	if err != nil && err != sql.ErrNoRows {
		log.Printf("something went wrong %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"messages": err.Error()})
	} else if err != nil && err == sql.ErrNoRows {
		log.Printf("no data found for query %s", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "min max values has not generated for the service"})
	} else {
		c.JSON(http.StatusOK, gin.H{"min": minVal, "max": maxVal})
	}
}

// GetLDMSummary api to get ldm service summary
func GetLDMSummary(c *gin.Context) {
	service := c.MustGet("service").(string)
	uuid := c.Param("uuid")
	ts := c.Param("timestamp")
	dbtype := c.Param("type")
	lyrtype := c.Param("lyrtype")

	var (
		err    error
		minVal float64
		maxVal float64
	)

	minVal, maxVal, err = model.GetSummary(service, uuid, ts, dbtype, lyrtype)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("something went wrong %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"messages": err.Error()})
	} else if err != nil && err == sql.ErrNoRows {
		log.Printf("no data found for query %s", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "min max values has not generated for the service"})
	} else {
		c.JSON(http.StatusOK, gin.H{"min": minVal, "max": maxVal})
	}
}

// GetSummaries returns summaries filtered by service, uuid and timestsamp
func GetSummaries(c *gin.Context) {
	uuids := c.QueryArray("uuid")
	timestampFrom := c.Query("timestamp_from")
	timestampUntil := c.Query("timestamp_until")
	removeZero := c.Query("remove_zero")
	var services []string

	services = c.QueryArray("service")
	if len(services) == 0 {
		services = []string{c.MustGet("service").(string)}
	}

	// Filter out any UUIDs that user don't have access to.
	// Need to do this in combination with each service.
	// We'll essentially trim the list if UUIDs for each service.
	for _, service := range services {
		uuids = lib.FilterUUIDsByService(c, service, uuids)
	}

	res, err := model.GetSummaries(services, uuids, timestampFrom, timestampUntil, removeZero != "")
	if err != nil && err != sql.ErrNoRows {
		log.Printf("something went wrong: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"messages": err.Error()})
	} else if err != nil && err == sql.ErrNoRows {
		log.Printf("no data found for query: %s", err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "min max values has not generated for the service"})
	} else {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"data": res})
	}
}
