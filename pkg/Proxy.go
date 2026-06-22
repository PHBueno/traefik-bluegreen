package pkg

import (
	"net/http/httputil"
	"net/url"

	"github.com/PHBueno/traefik-bluegreen/pkg/redis"
)

type Proxy struct {
	ProxyURL  *url.URL
	RedisConn *redis.RedisStore
}

func (p *Proxy) RewriteProxy() func(*httputil.ProxyRequest) {

	return func(pr *httputil.ProxyRequest) {
		pr.SetURL(p.ProxyURL)
		pr.Out.Host = pr.In.Host

		tenant := pr.In.URL.Query().Get("tenant")
		app := pr.In.Header.Get("X-App-Slug")

		tenantModel := p.RedisConn.GetSlot(tenant, app)

		pr.Out.Header.Set("X-Slot", tenantModel.Slot)

		pr.SetXForwarded()
	}
}
