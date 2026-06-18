package redis

import (
	"fmt"
	"net"
	"os"
	"sync"
)

var (
	store *RedisStore
	once  sync.Once
)

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
	once.Do(
		func() {
			fmt.Fprintln(os.Stdout, "[ONCE] => Executando Once")
			store = &RedisStore{
				address: address,
				port:    port,
				cache:   make(map[string]map[string]string),
			}
		},
	)
	return store

}

func (rs *RedisStore) GetSlot(tenant string, app string) *TenantSlot {
	fmt.Fprintln(os.Stdout, rs.cache)
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

	fmt.Fprintln(os.Stdout, "[REDIS CACHE] => valor encontrado no cache")

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
	rs.updateCache(tenant, app, "1")

	return &TenantSlot{
		TenantID: "ID",
		AppName:  "appName",
		Slot:     "1",
	}, nil

}

func (rs *RedisStore) updateCache(tenant string, app string, slot string) {
	rs.cache[fmt.Sprintf("%s-%s", tenant, app)] = map[string]string{
		"tenantID": tenant,
		"appName":  app,
		"slot":     slot,
	}
	fmt.Fprintln(os.Stdout, "[REDIS CACHE] => cache atualizado com sucesso!")
}
