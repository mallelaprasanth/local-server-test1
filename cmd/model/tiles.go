package model

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/local-server-test1/database"
)

// get tile data from postgres
func GetTiles(service string, uuid string, timestamp string, z uint64, x uint64, y uint64) (data []byte, err error) {
	var rds = database.GetRedisConn()
	tck := fmt.Sprintf("%d/%d/%d", z, x, y)
	layer := fmt.Sprintf("%s_%s_%s", service, uuid, timestamp)

	// Check cache
	hgetRes := rds.HGet(layer, tck)

	hit := hgetRes.Err() != redis.Nil && hgetRes.Err() == nil
	if hgetRes.Err() != redis.Nil && hgetRes.Err() != nil {
		// Don't error out but log it
		log.Printf("Failed to HGET cached data for layer %s, tile %s: %s", layer, tck, hgetRes.Err())
	}

	if hit {
		data, err = hgetRes.Bytes()
		if len(data) > 0 {
			println("serving from redis, -------------->")
		}
		// Don't error out but log it
		if err != nil {
			log.Printf("Failed to HGET cached data for layer %s, tile %s: %s", layer, tck, hgetRes.Err())
		}
	}
	if len(data) == 0 {
		var db = database.GetConnection()
		err = db.Ping()
		if err != nil {
			return data, err
		}
		y2 := (1 << z) - 1 - y
		//prepare
		sqlSelect := "select tile_data from tiles where service =? and uuid=? and timestamp=? and zoom_level = ? and tile_column = ? and tile_row = ?"
		row := db.QueryRow(sqlSelect, service, uuid, timestamp, z, x, y2)
		err = row.Scan(&data)

		if len(data) != 0 {
			// Cache data
			hsetRes := rds.HSet(layer, tck, data)
			println("cache it in redis, ----------------------->")
			if hsetRes.Err() != redis.Nil && hsetRes.Err() != nil {
				// Don't error out but log it
				log.Printf("Failed to HSET cache data for layer %s, tile %s: %s", layer, tck, hsetRes.Err())
			}
		}
	}
	return data, err
}

// get summary data for service
func GetSummary(tilefolder string, service string, uuid string, timestamp string) (data []byte, err error) {
	data, err = ioutil.ReadFile(os.Getenv("TILES_FOLDER") + "/" + string(service) + "/" + string(uuid) + "/minmax.json")
	return data, err
}
