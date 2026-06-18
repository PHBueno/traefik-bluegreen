package pkg

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"

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

		slot := p.RedisConn.GetSlot(tenant, app)

		switch tenant {
		case "456":
			pr.Out.Header.Set("X-Slot", slot.Slot)
		default:
			pr.Out.Header.Set("X-Slot", slot.Slot)
		}

		pr.SetXForwarded()

		fmt.Fprintf(os.Stdout,
			"Encaminhando requisição para o Traefik -> Host: %s | Headers: %s\n",
			pr.Out.Host, pr.Out.Header,
		)
	}
}
