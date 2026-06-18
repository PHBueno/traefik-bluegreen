package redis

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/PHBueno/traefik-bluegreen/pkg/redis/models"
)

var (
	store *RedisStore
	once  sync.Once
)

type RedisStore struct {
	address string
	port    string
	cache   map[string]*models.TenantSlot
	mu      sync.RWMutex
}

func NewConnection(address string, port string) *RedisStore {
	once.Do(
		func() {
			store = &RedisStore{
				address: address,
				port:    port,
				cache:   make(map[string]*models.TenantSlot),
			}
		},
	)
	return store

}

func (rs *RedisStore) GetSlot(tenant string, app string) *models.TenantSlot {
	fmt.Fprintln(os.Stdout, rs.cache)
	// tenta buscar do cache
	tenantSlot, err := rs.getCachedSlot(tenant, app)

	if err != nil {
		// se não existir cache, busca do redis
		tenantSlot, _ = rs.getRedisSlot(tenant, app)
	}

	return tenantSlot

}

// Busca valores do Cache
func (rs *RedisStore) getCachedSlot(tenant string, app string) (*models.TenantSlot, error) {
	rs.mu.RLock() // Protege a leitura do cache permitindo multiplas leituras
	tenantData, tenantExists := rs.cache[fmt.Sprintf("%s:%s", tenant, app)]
	rs.mu.RUnlock()

	if !tenantExists {
		fmt.Fprintln(os.Stderr, "[REDIS CACHE] => valor não encontrado no cache")
		return nil, fmt.Errorf("valor não encontrado no cache!")
	}

	fmt.Fprintln(os.Stdout, "[REDIS CACHE] => valor encontrado no cache")

	return &models.TenantSlot{
		TenantID: tenantData.TenantID,
		AppName:  tenantData.AppName,
		Slot:     tenantData.Slot,
	}, nil
}

// Busca valores do Redis
func (rs *RedisStore) getRedisSlot(tenant string, app string) (*models.TenantSlot, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(rs.address, rs.port))

	if err != nil {
		fmt.Fprintln(os.Stderr, "[REDIS CONNECTION] => erro para conectar ao redis: ", err)
		return nil, err
	}

	defer conn.Close() // fecha a conexão após o retorno da função.
	fmt.Fprintln(os.Stdout, "[REDIS CONNECTION] => conexão estabelecida com sucesso")
	HGetAll(conn, fmt.Sprintf("%s:%s", tenant, app))

	tenantModel, _ := Deserializer(conn)

	rs.updateCache(tenant, app, tenantModel.Slot)

	return tenantModel, nil
}

// Atualiza Cache
func (rs *RedisStore) updateCache(tenant string, app string, slot string) {
	rs.mu.Lock() // Protege escrita do cache
	rs.cache[fmt.Sprintf("%s:%s", tenant, app)] = &models.TenantSlot{
		TenantID: tenant,
		AppName:  app,
		Slot:     slot,
	}
	rs.mu.Unlock()

	fmt.Fprintln(os.Stdout, "[REDIS CACHE] => cache atualizado com sucesso!")
}
