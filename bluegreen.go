package traefik_bluegreen

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/PHBueno/traefik-bluegreen/pkg"
)

type Config struct {
	RedisAddress  string
	RedisPort     string
	RedisPassword string
	RedisDataBase string
}

func CreateConfig() *Config {
	return &Config{}
}

func rewriteProxy(traefikTarget *url.URL) func(*httputil.ProxyRequest) {
	return func(pr *httputil.ProxyRequest) {
		pr.SetURL(traefikTarget)
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

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.RedisAddress == "" {
		return nil, fmt.Errorf("Redis Address is not set!")
	}

	traefikTarget, err := url.Parse("https://traefik.traefik-controller.svc.cluster.local:443")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	proxy := &httputil.ReverseProxy{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Rewrite: rewriteProxy(traefikTarget),
	}

	bg := pkg.New(next, proxy, name)

	return bg, nil
}
