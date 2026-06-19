package commands

import (
	"fmt"
	"net"
	"os"

	"github.com/PHBueno/traefik-bluegreen/pkg/redis/models"
	resp "github.com/PHBueno/traefik-bluegreen/pkg/redis/resp"
)

func HGetAll(conn net.Conn, key string) (*models.TenantSlot, error) {

	_, err := conn.Write([]byte(resp.Serializer("HGETALL", key)))

	if err != nil {
		fmt.Fprintln(os.Stdout, "[REDIS CONNECTION] => erro para buscar valor no redis: ", err)
		return nil, err
	}

	tenantModel, err := resp.Deserializer(conn)

	if err != nil {
		fmt.Fprintln(os.Stdout, "[REDIS CONNECTION] => erro na resposta do redis: ", err)
		return nil, err
	}

	return tenantModel, nil

}
