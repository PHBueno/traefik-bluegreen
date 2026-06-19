package redis

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/PHBueno/traefik-bluegreen/pkg/redis/cache"
	"github.com/PHBueno/traefik-bluegreen/pkg/redis/commands"
	"github.com/PHBueno/traefik-bluegreen/pkg/redis/models"
)

type RedisStore struct {
	address    string
	port       string
	localCache *cache.LocalCache // TODO: Adicionar invalidação do Cache
}

func NewConnection(address string, port string) *RedisStore {
	return &RedisStore{
		address:    address,
		port:       port,
		localCache: cache.NewLocalCache(),
	}
}

func (rs *RedisStore) GetSlot(tenant string, app string) *models.TenantSlot {
	fmt.Fprintln(os.Stdout, rs.localCache)
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
	conn, err := net.Dial("tcp", net.JoinHostPort(rs.address, rs.port))

	if err != nil {
		slog.Error("[REDIS CONNECTION] => erro para conectar ao redis", "error", err)
		return nil, err
	}

	defer conn.Close() // fecha a conexão após o retorno da função.

	slog.Info("[REDIS CONNECTION] => conexão estabelecida com sucesso")

	tenantModel, _ := commands.HGetAll(conn, fmt.Sprintf("%s:%s", tenant, app))

	rs.updateCache(tenant, app, tenantModel.Slot)

	return tenantModel, nil
}

// Atualiza Cache
func (rs *RedisStore) updateCache(tenant string, app string, slot string) {
	rs.localCache.SetTenant(&models.TenantSlot{
		TenantID: tenant,
		AppName:  app,
		Slot:     slot,
	})

	slog.Info("[REDIS CACHE] => cache atualizado com sucesso!")
}
