package redis

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func HGetAll(conn net.Conn, key string) error {
	var stringBuilder strings.Builder

	stringBuilder.WriteString(fmt.Sprintf("*2\r\n$7\r\nHGETALL\r\n$%d\r\n%s\r\n", len(key), key))

	_, err := conn.Write([]byte(stringBuilder.String()))

	if err != nil {
		fmt.Fprintln(os.Stdout, "[REDIS CONNECTION] => erro para escrever no redis: ", err)
		return err
	}

	fmt.Fprintln(os.Stdout, "[REDIS CONNECTION] => sucesso para escrever no redis")
	return nil

}
