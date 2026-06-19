package redis

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func Serializer(command string, args ...string) string {
	var stringBuilder strings.Builder

	fmt.Fprintf(
		&stringBuilder, "*%d\r\n$%d\r\n%s\r\n",
		len(args)+1, len(command), command,
	) // command + quantidade de args

	for i := 0; i < len(args); i++ {
		fmt.Fprintf(
			&stringBuilder, "$%d\r\n%s\r\n",
			len(args[i]), args[i],
		)
	}

	return stringBuilder.String()

}

func HGetAll(conn net.Conn, key string) error {
	var stringBuilder strings.Builder

	fmt.Fprintf(
		&stringBuilder,
		"*2\r\n$7\r\nHGETALL\r\n$%d\r\n%s\r\n",
		len(key),
		key,
	)

	_, err := conn.Write([]byte(stringBuilder.String()))

	if err != nil {
		fmt.Fprintln(os.Stdout, "[REDIS CONNECTION] => erro para escrever no redis: ", err)
		return err
	}

	fmt.Fprintln(os.Stdout, "[REDIS CONNECTION] => sucesso para escrever no redis")
	return nil

}
