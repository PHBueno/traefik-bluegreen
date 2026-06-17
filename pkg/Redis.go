package pkg

import (
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
)

// type ClientTenant struct {
// 	app    string
// 	slot   string
// 	tenant string
// }

func RedisConn(address string, password string, database int) (redis.Conn, error) {
	conn, err := redis.Dial(
		"tcp",
		address,
		redis.DialPassword(password),
		redis.DialDatabase(database),
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, "[REDIS ERROR] => ", err)
		return nil, err
	}

	data, err := conn.Do("PING")

	if err != nil {
		fmt.Fprintln(os.Stderr, "[REDIS ERROR] => ", err)
		return nil, err
	}

	fmt.Fprintln(os.Stdout, "[REDIS CONNECTION] => ", data)

	return conn, nil
}

// func NewRedisConnection(address string, password string, database int) (*redis.Client, error) {
// 	context := context.Background()

// 	redisClient := redis.NewClient(
// 		&redis.Options{
// 			Addr:     address,
// 			Password: password,
// 			DB:       database,
// 		},
// 	)

// 	_, err := redisClient.Ping(context).Result()

// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, "[REDIS ERROR] => ", err)
// 		return nil, err
// 	}

// 	fmt.Fprintln(os.Stdout, "Successfully connecting to Redis!")
// 	return redisClient, nil
// }
