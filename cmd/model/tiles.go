package model

import (
	"bytes"
	"compress/zlib"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"github.com/lib/pq"
	"github.com/local-server-test1/cmd/database"
)

// get tile data from postgres
func GetTiles(service string, uuid string, timestamp string, z uint64, x uint64, y uint64, subtype string, lyrtype string) (data []byte, err error) {
	var tmpByteArr []byte
	var compressed bool
	var layer string

	var dbtype string
	var zoomlevel string

	var rds = database.GetRedisConn()

	tck := fmt.Sprintf("%d/%d/%d", z, x, y)
	layer = fmt.Sprintf("%s_%s_%s", service, uuid, timestamp)
	if subtype != "" && lyrtype != "" {
		layer = fmt.Sprintf("%s_%s_%s_%s_%s_%d", service, uuid, subtype, timestamp, lyrtype, z)
		dbtype = subtype + "-" + lyrtype
		zoomlevel = "-" + strconv.FormatUint(z, 10)
	}

	// Check cache
	hgetRes := rds.HGet(layer, tck)

	hit := hgetRes.Err() != redis.Nil && hgetRes.Err() == nil
	if hgetRes.Err() != redis.Nil && hgetRes.Err() != nil {
		// Don't error out but log it
		log.Printf("Failed to HGET cached data for layer %s, tile %s: %s", layer, tck, hgetRes.Err())
	}

	if hit {
		tmpByteArr, err = hgetRes.Bytes()
		if len(tmpByteArr) > 0 {
			data = tmpByteArr
			log.Printf("serving from redis, -------------->")
		}
		// Don't error out but log it
		if err != nil {
			log.Printf("Failed to HGET cached data for layer %s, tile %s: %s", layer, tck, hgetRes.Err())
		}
	}

	if len(tmpByteArr) == 0 {
		var db = database.GetConnection()
		err = db.Ping()
		if err != nil {
			return tmpByteArr, err
		}
		y2 := (1 << z) - 1 - y
		//prepare
		var sqlSelect string
		var row *sql.Row

		if subtype != "" {
			if strings.Contains(dbtype, "markers") {
				sqlSelect = "select tile_data,COALESCE(compressed,False) as compressed from tiles where service=$1 and uuid=$2 and timestamp=$3 and zoom_level=$4 and tile_column=$5 and tile_row=$6 and type=$7"
				row = db.QueryRow(sqlSelect, service, uuid, timestamp, z, x, y2, dbtype)
			} else {
				sqlSelect = "select tile_data,COALESCE(compressed,False) as compressed from tiles where service=$1 and uuid=$2 and timestamp=$3 and zoom_level=$4 and tile_column=$5 and tile_row=$6 and type ilike $7 || '%' || $8 || '%'"
				row = db.QueryRow(sqlSelect, service, uuid, timestamp, z, x, y2, dbtype, zoomlevel)
			}
		} else {
			sqlSelect = "select tile_data,COALESCE(compressed,False) as compressed from tiles where service=$1 and uuid=$2 and timestamp=$3 and zoom_level=$4 and tile_column=$5 and tile_row=$6"
			row = db.QueryRow(sqlSelect, service, uuid, timestamp, z, x, y2)
		}

		err = row.Scan(&tmpByteArr, &compressed)

		if err != nil && err == sql.ErrNoRows {
			log.Printf("no data found for query, select tile_data from tiles where service ='%s' and uuid='%s' and timestamp='%s' and zoom_level = %d and tile_column = %d and tile_row = %d and type ilike '%s'", service, uuid, timestamp, z, x, y2, dbtype+zoomlevel)
			log.Printf(err.Error())
		} else {
			data = tmpByteArr

			// uncompress the data
			if compressed {

				b := bytes.NewReader(tmpByteArr)
				r, err := zlib.NewReader(b)
				if err != nil {
					panic(err)
				}
				defer r.Close()

				p, err := ioutil.ReadAll(r)
				if err != nil {
					panic(err)
				}
				data = p
			}

			// Cache data
			hsetRes := rds.HSet(layer, tck, data)
			log.Printf("cache it in redis, ----------------------->")
			if hsetRes.Err() != redis.Nil && hsetRes.Err() != nil {
				// Don't error out but log it
				log.Printf("Failed to HSET cache data for layer %s, tile %s: %s", layer, tck, hsetRes.Err())
			}
		}
	}
	return data, err
}

// get summary data for service
func GetSummary(service string, uuid string, timestamp string, subtype string, lyrtype string) (minVal float64, maxVal float64, err error) {
	var dbtype string
	var db = database.GetConnection()
	err = db.Ping()
	if err != nil {
		return minVal, maxVal, err
	}

	if subtype != "" && lyrtype != "" {
		dbtype = subtype + "-" + lyrtype
	}

	//prepare
	var sqlSelect string
	var row *sql.Row

	if subtype != "" {
		sqlSelect = "select min_val,max_val from summary where service=$1 and uuid=$2 and timestamp=$3 and type=$4"
		row = db.QueryRow(sqlSelect, service, uuid, timestamp, dbtype)
	} else {
		sqlSelect = "select min_val,max_val from summary where service=$1 and uuid=$2 and timestamp=$3"
		row = db.QueryRow(sqlSelect, service, uuid, timestamp)
	}

	err = row.Scan(&minVal, &maxVal)

	if err != nil || err == sql.ErrNoRows {
		log.Printf("no data found for query, select min_val,max_val from summary where service='%s' and uuid='%s' and timestamp='%s' and type='%s'", service, uuid, timestamp, dbtype)
	}

	return minVal, maxVal, err
}

// GetSummaries returns summaries filtered by service, uuid and timestsamp
func GetSummaries(services []string, uuids []string, timestampFrom string, timestampUntil string, removeZero bool) (res []map[string]interface{}, err error) {
	var db = database.GetConnection()
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	//prepare
	res = []map[string]interface{}{}
	sqlSelect := "SELECT service, uuid, timestamp, min_val, max_val FROM summary WHERE service = ANY($1) AND uuid = ANY($2) AND timestamp BETWEEN $3 AND $4"
	rows, err := db.Query(sqlSelect, pq.Array(services), pq.Array(uuids), timestampFrom, timestampUntil)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			service   string
			uuid      uint64
			timestamp string
			min       float64
			max       float64
		)

		if err := rows.Scan(&service, &uuid, &timestamp, &min, &max); err != nil {
			return nil, err
		}

		d := map[string]interface{}{
			"service":   service,
			"uuid":      uuid,
			"timestamp": timestamp,
			"min":       min,
			"max":       max,
		}

		if removeZero {
			if isNotZero(service, min, max) {
				res = append(res, d)
			}
		} else {
			res = append(res, d)
		}
	}

	return res, err
}

func isNotZero(service string, min, max float64) bool {
	return min != max
}
