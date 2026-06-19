package traefik_bluegreen

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/PHBueno/traefik-bluegreen/pkg"
	"github.com/PHBueno/traefik-bluegreen/pkg/redis"
)

type Config struct {
	RedisAddress  string
	RedisPort     string
	RedisPassword string
	RedisDataBase string
}

func CreateConfig() *Config {
	return &Config{
		RedisPort: "6379",
	}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if config.RedisAddress == "" {
		slog.Error("[REDIS CONFIG] The Redis address has not been set")
		return nil, fmt.Errorf("[REDIS CONFIG] The Redis address has not been set")
	}

	traefikTarget, err := url.Parse("https://traefik.traefik-controller.svc.cluster.local:443")

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	redisConn := redis.NewConnection(config.RedisAddress, config.RedisPort)

	targetProxy := &pkg.Proxy{
		ProxyURL:  traefikTarget,
		RedisConn: redisConn,
	}

	proxy := &httputil.ReverseProxy{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Rewrite: targetProxy.RewriteProxy(),
	}

	bg := pkg.New(next, proxy, name)

	return bg, nil
}
