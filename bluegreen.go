package traefik_bluegreen

import (
	"context"
	"crypto/tls"
	"fmt"
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

func init() {
	fmt.Println("PLUGIN INIT")
}

func CreateConfig() *Config {
	return &Config{}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.RedisAddress == "" {
		return nil, fmt.Errorf("Redis Address is not set!")
	}

	fmt.Fprintf(os.Stdout, "NAME => %s", name)

	traefikTarget, err := url.Parse("https://traefik.traefik-controller.svc.cluster.local:443")

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	redisConn, err := redis.NewConnection(config.RedisAddress, config.RedisPort)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

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

// Se a requisição não vier com app-id, deve ser encaminhado para o um Default;
// Se não vier tenant, tem que buscar pela app Default;
// definir uma espécie de Default Backend;

// Cenários onde não vier o app-id;
// Cenários onde não vier o tenant;
// Cenários onde não vier nem o app-id e nem o tenant-id;
