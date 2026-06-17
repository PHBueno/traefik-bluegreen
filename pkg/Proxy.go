package pkg

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/redis/go-redis/v9"
)

type Traefik struct {
	URL *url.URL
}

func (t *Traefik) RewriteProxy(redisConn *redis.Client) func(*httputil.ProxyRequest) {

	return func(pr *httputil.ProxyRequest) {
		pr.SetURL(t.URL)
		pr.Out.Host = pr.In.Host

		tenant := pr.In.URL.Query().Get("tenant")

		switch tenant {
		case "456":
			pr.Out.Header.Set("X-Slot", "1")
		default:
			pr.Out.Header.Set("X-Slot", "2")
		}

		pr.SetXForwarded()

		fmt.Fprintf(os.Stdout,
			"Encaminhando requisição para o Traefik -> Host: %s | Headers: %s\n",
			pr.Out.Host, pr.Out.Header,
		)
	}
}
