package pkg

import (
	"log/slog"
	"net/http/httputil"
	"net/url"

	"github.com/PHBueno/traefik-bluegreen/pkg/redis"
)

type Proxy struct {
	ProxyURL  *url.URL
	RedisConn *redis.RedisStore
}

func verifyEmpty(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func (p *Proxy) RewriteProxy() func(*httputil.ProxyRequest) {

	return func(pr *httputil.ProxyRequest) {
		pr.SetURL(p.ProxyURL)
		pr.Out.Host = pr.In.Host

		tenant := verifyEmpty(pr.In.URL.Query().Get("tenant"), "000000") // tenant default => 000000
		app := verifyEmpty(pr.In.Header.Get("X-App-Slug"), "default")    // app default => default

		slog.Info("Application Path", "path", pr.In.URL.Path)

		tenantModel, err := p.RedisConn.GetSlot(tenant, app)
		slot := "-1" // Caso não encontre o valor no Redis nem no Cache

		if err != nil {
			slog.Error(err.Error())
		} else {
			slot = tenantModel.Slot
		}

		pr.Out.Header.Set("X-Slot", slot)

		pr.SetXForwarded()
	}
}
