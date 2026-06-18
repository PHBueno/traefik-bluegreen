package redis

import (
	"fmt"
	"net"
	"os"
)

// Conexão com o Redis

type RedisStore struct {
	address string
	port    string
	cache   map[string]map[string]string
}

/*
	123-nginx: {
	  tenantID: 123
	  appName: nginx
	  slot: 1
	}
*/
type TenantSlot struct {
	TenantID string
	AppName  string
	Slot     string
}

func NewConnection(address string, port string) *RedisStore {
	return &RedisStore{
		address: address,
		port:    port,
		cache:   make(map[string]map[string]string),
	}

}

func (rs *RedisStore) GetSlot(tenant string, app string) *TenantSlot {
	tenantSlot, err := rs.getCachedSlot(tenant, app)

	if err != nil {
		tenantSlot, _ = rs.getRedisSlot(tenant, app)
	}

	return tenantSlot

}

func (rs *RedisStore) getCachedSlot(tenant string, app string) (*TenantSlot, error) {
	tenantData, tenantExists := rs.cache[fmt.Sprintf("%s-%s", tenant, app)]

	if !tenantExists {
		fmt.Fprintln(os.Stderr, "[REDIS CACHE] => valor não encontrado no cache")
		return nil, fmt.Errorf("valor não encontrado no cache!")
	}

	return &TenantSlot{
		TenantID: tenantData["tenantID"],
		AppName:  tenantData["appName"],
		Slot:     tenantData["slot"],
	}, nil
}

func (rs *RedisStore) getRedisSlot(tenant string, app string) (*TenantSlot, error) {
	_, err := net.Dial("tcp", net.JoinHostPort(rs.address, rs.port))

	if err != nil {
		fmt.Fprintln(os.Stderr, "[REDIS CONNECTION] => erro para conectar ao redis: ", err)
		return nil, err
	}

	fmt.Fprintln(os.Stdout, "[REDIS CONNECTION] => conexão estabelecida com sucesso")
	return &TenantSlot{
		TenantID: "ID",
		AppName:  "appName",
		Slot:     "1",
	}, nil

}
