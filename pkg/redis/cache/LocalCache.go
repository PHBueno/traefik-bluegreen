package cache

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/PHBueno/traefik-bluegreen/pkg/redis/models"
)

var (
	cache *LocalCache
	once  sync.Once
)

type cacheEntry struct {
	tenant    *models.TenantSlot
	expiresAt time.Time
}

type LocalCache struct {
	cache map[string]*cacheEntry
	mu    sync.RWMutex
}

func NewLocalCache() *LocalCache {
	once.Do(
		func() {
			cache = &LocalCache{
				cache: make(map[string]*cacheEntry),
			}
		},
	)
	return cache
}

// Escrita
func (lc *LocalCache) SetTenant(tenant *models.TenantSlot, ttl int) {
	lc.mu.Lock()

	lc.cache[fmt.Sprintf("%s:%s", tenant.TenantID, tenant.AppName)] = &cacheEntry{
		tenant:    tenant,
		expiresAt: time.Now().Add(time.Duration(ttl) * time.Second),
	}

	lc.mu.Unlock()
}

// Leitura
func (lc *LocalCache) GetTenant(Id string) (*models.TenantSlot, error) {
	lc.mu.RLock()
	entry, exists := lc.cache[Id]

	if !exists {
		lc.mu.RUnlock()
		slog.Info("[REDIS CACHE] => valor não encontrado no cache")
		return nil, fmt.Errorf("valor não encontrado no cache!")
	}

	if time.Now().After(entry.expiresAt) {
		lc.mu.RUnlock()

		lc.mu.Lock()

		// Garante que não teve modificações concorrentes entre uma escrita e outra
		if current, ok := lc.cache[Id]; ok && current == entry {
			delete(lc.cache, Id)
		}

		lc.mu.Unlock()
		slog.Info("[REDIS CACHE] => cache inválido!")
		return nil, fmt.Errorf("cache inválido!")

	}

	tenant := entry.tenant
	lc.mu.RUnlock()

	return &models.TenantSlot{
		TenantID: tenant.TenantID,
		AppName:  tenant.AppName,
		Slot:     tenant.Slot,
	}, nil
}
