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
	cache   map[string]*TenantSlot
	mu      sync.RWMutex
}

type TenantSlot struct {
	TenantID string
	AppName  string
	Slot     string
}

func NewConnection(address string, port string) *RedisStore {
	once.Do(
		func() {
			store = &RedisStore{
				address: address,
				port:    port,
				cache:   make(map[string]*TenantSlot),
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
	rs.mu.RLock() // Protege a leitura do cache permitindo multiplas leituras
	tenantData, tenantExists := rs.cache[fmt.Sprintf("%s:%s", tenant, app)]
	rs.mu.RUnlock()

	if !tenantExists {
		fmt.Fprintln(os.Stderr, "[REDIS CACHE] => valor não encontrado no cache")
		return nil, fmt.Errorf("valor não encontrado no cache!")
	}

	fmt.Fprintln(os.Stdout, "[REDIS CACHE] => valor encontrado no cache")

	return &TenantSlot{
		TenantID: tenantData.TenantID,
		AppName:  tenantData.AppName,
		Slot:     tenantData.Slot,
	}, nil
}

func (rs *RedisStore) getRedisSlot(tenant string, app string) (*TenantSlot, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(rs.address, rs.port))

	if err != nil {
		fmt.Fprintln(os.Stderr, "[REDIS CONNECTION] => erro para conectar ao redis: ", err)
		return nil, err
	}

	defer conn.Close() // fecha a conexão após o retorno da função.

	HGetAll(conn, fmt.Sprintf("%s:%s", tenant, app))

	fmt.Fprintln(os.Stdout, "[REDIS CONNECTION] => conexão estabelecida com sucesso")
	rs.updateCache(tenant, app, "1")

	return &TenantSlot{
		TenantID: "ID",
		AppName:  "appName",
		Slot:     "1",
	}, nil

}

func (rs *RedisStore) updateCache(tenant string, app string, slot string) {
	rs.mu.Lock() // Protege escrita do cache
	rs.cache[fmt.Sprintf("%s:%s", tenant, app)] = &TenantSlot{
		TenantID: tenant,
		AppName:  app,
		Slot:     slot,
	}
	rs.mu.Unlock()

	fmt.Fprintln(os.Stdout, "[REDIS CACHE] => cache atualizado com sucesso!")
}
