package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

const (
	retryInterval = 5
	backOffLimit  = 4
)

var (
	conn *sql.DB
	rds  *redis.Client
	err  error
)

//connect to postgres
func Connect(psqlInfo string) (err error) {
	var retries = 0
	for i := 0; i < backOffLimit; i++ {
		conn, err = sql.Open("postgres", psqlInfo)

		if err == nil {
			break
		}

		wait := math.Pow(2, float64(retries)) * retryInterval
		log.Printf("Failed to connect to DB. Wait for %.0fs. Retry #%d.", wait, retries)
		time.Sleep(time.Duration(wait) * time.Second)
		retries++
	}

	if err != nil {
		return err
	}

	return nil
}

// Init is initialize db from main function
func Init() {
	// port, err := strconv.ParseUint(os.Getenv("PG_PORT"), 10, 32)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("PG_DBNAME"))

	err := Connect(psqlInfo)
	if err != nil {
		panic(err)
	}

	err = ConnectRedis()
	if err != nil {
		panic(err)
	}
}

func TestInit() {
	// port, err := strconv.ParseUint(os.Getenv("PG_PORT"), 10, 32)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("PG_DBNAME"))

	err := Connect(psqlInfo)
	if err != nil {
		panic(err)
	}

	err = ConnectRedis()
	if err != nil {
		panic(err)
	}
}

// GetDB is called in models
func GetConnection() *sql.DB {
	return conn
}

// Close is closing db
func Close() {
	if err := conn.Close(); err != nil {
		panic(err)
	}
}

// GetDB is called in models
func GetRedisConn() *redis.Client {
	return rds
}

// connect to redis
func ConnectRedis() error {
	rdsDB, err1 := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 32)
	if err1 != nil {
		return errors.New("Failed to parse REDIS config")
	}

	rdsPrt := os.Getenv("REDIS_PORT")
	addr := os.Getenv("REDIS_HOST") + ":" + rdsPrt

	log.Printf("Connecting to Redis at %s using DB %d", addr, rdsDB)

	rds = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASS"),
		DB:       int(rdsDB),
	})

	log.Printf("connected to Redis using DB ")
	return nil
}
