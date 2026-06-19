package cache

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/PHBueno/traefik-bluegreen/pkg/redis/models"
)

var (
	cache *LocalCache
	once  sync.Once
)

type LocalCache struct {
	cache map[string]*models.TenantSlot
	mu    sync.RWMutex
}

func NewLocalCache() *LocalCache {
	once.Do(
		func() {
			cache = &LocalCache{
				cache: make(map[string]*models.TenantSlot),
			}
		},
	)
	return cache
}

func (lc *LocalCache) SetTenant(tenant *models.TenantSlot) {
	lc.mu.RLock()
	lc.cache[fmt.Sprintf("%s:%s", tenant.TenantID, tenant.AppName)] = tenant
	lc.mu.RUnlock()
}

func (lc *LocalCache) GetTenant(Id string) (*models.TenantSlot, error) {
	lc.mu.RLock()
	tenant, exists := lc.cache[Id]
	lc.mu.RUnlock()

	if !exists {
		slog.Info("[REDIS CACHE] => valor não encontrado no cache")
		return nil, fmt.Errorf("valor não encontrado no cache!")
	}

	return tenant, nil
}
