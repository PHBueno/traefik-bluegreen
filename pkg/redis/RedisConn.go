package redis

import (
	"fmt"
	"net"
	"os"
)

// Conexão com o Redis

type RedisConn struct {
	connection net.Conn
}

func NewConnection(address string, port string) (*RedisConn, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(address, port))

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	fmt.Fprintln(os.Stdout, "connection established with Redis")
	return &RedisConn{connection: conn}, nil

}
