package redis

import (
	"fmt"
	"net"
	"os"
)

// Conexão com o Redis

type RedisConn struct {
	Address string
	Port    string
}

func (r *RedisConn) NewConnection() (*net.Conn, error) {
	connectionString := r.Address + ":" + r.Port

	conn, err := net.Dial("tcp", connectionString)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	fmt.Fprintln(os.Stdout, "connection established with Redis")
	return &conn, nil
}
