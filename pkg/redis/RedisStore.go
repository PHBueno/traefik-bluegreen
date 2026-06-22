package redis

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/PHBueno/traefik-bluegreen/pkg/redis/cache"
	"github.com/PHBueno/traefik-bluegreen/pkg/redis/commands"
	"github.com/PHBueno/traefik-bluegreen/pkg/redis/models"
)

type RedisStore struct {
	address    string
	port       string
	localCache *cache.LocalCache
	cacheTTL   int
}

func NewConnection(address string, port string, cacheTTL int) *RedisStore {
	return &RedisStore{
		address:    address,
		port:       port,
		localCache: cache.NewLocalCache(),
		cacheTTL:   cacheTTL,
	}
}

func verifyEmpty(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func (rs *RedisStore) GetSlot(tenant string, app string) (*models.TenantSlot, error) {

	tenant = verifyEmpty(tenant, "000000")
	app = verifyEmpty(app, "default")

	// tenta buscar do cache
	tenantSlot, err := rs.getCachedSlot(tenant, app)

	if err != nil {
		// se não existir cache, busca do redis
		tenantSlot, err = rs.getRedisSlot(tenant, app)

		if err != nil {
			return nil, err
		}
	}

	return tenantSlot, nil

}

// Busca valores do Cache
func (rs *RedisStore) getCachedSlot(tenant string, app string) (*models.TenantSlot, error) {
	tenantData, err := rs.localCache.GetTenant(fmt.Sprintf("%s:%s", tenant, app))

	if err != nil {
		return nil, fmt.Errorf("[REDIS CACHE] => valor não encontrado no cache!")
	}

	slog.Info("[REDIS CACHE] => valor encontrado no cache")

	return &models.TenantSlot{
		TenantID: tenantData.TenantID,
		AppName:  tenantData.AppName,
		Slot:     tenantData.Slot,
	}, nil
}

// Busca valores do Redis
func (rs *RedisStore) getRedisSlot(tenant string, app string) (*models.TenantSlot, error) {
	conn, err := rs.redisConn()

	if err != nil {
		return nil, err
	}

	slog.Info("[REDIS CONNECTION] => conexão estabelecida com sucesso")

	tenantModel, err := commands.HGetAll(conn, fmt.Sprintf("%s:%s", tenant, app))
	conn.Close() // fecha a conexão após o retorno da função.
	slog.Info("[REDIS CONNECTION] => conexão fechada com sucesso")

	if err != nil {
		return nil, err
	}

	rs.updateCache(tenant, app, tenantModel.Slot)

	return tenantModel, nil
}

// Atualiza Cache
func (rs *RedisStore) updateCache(tenant string, app string, slot string) {
	rs.localCache.SetTenant(
		&models.TenantSlot{
			TenantID: tenant,
			AppName:  app,
			Slot:     slot,
		},
		rs.cacheTTL,
	)

	slog.Info("[REDIS CACHE] => cache atualizado com sucesso!")
}

func (rs *RedisStore) redisConn() (net.Conn, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(rs.address, rs.port))

	if err != nil {
		slog.Error("[REDIS CONNECTION] => erro para conectar ao redis", "error", err)
		return nil, err
	}

	return conn, nil
}
