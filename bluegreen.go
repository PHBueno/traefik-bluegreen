package traefik_bluegreen

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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

type BlueGreen struct {
	next  http.Handler
	proxy *httputil.ReverseProxy
	name  string
}

func (bg *BlueGreen) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(os.Stdout, "HEADERS ORIGINAIS: ", req.Header)
	fmt.Fprintln(os.Stdout, "TESTE => Chamando o ServeHTTP")
	bg.proxy.ServeHTTP(rw, req)
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	fmt.Fprintln(os.Stdout, "Iniciando contexto")
	if config.RedisAddress == "" {
		return nil, fmt.Errorf("Redis Address is not set!")
	}

	traefikTarget, err := url.Parse("http://traefik.traefik-controller.svc.cluster.local:80")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Println("Sucesso para acessar traefik")
	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.Out.Header.Del("X-Forwarded-Server")
			pr.Out.Header.Del("X-Forwarded-Host")
			pr.Out.Header.Del("X-Forwarded-Port")
			pr.Out.Header.Del("X-Forwarded-Proto")
			pr.Out.Header.Del("X-Forwarded-For")

			pr.SetURL(traefikTarget)
			pr.Out.Host = pr.In.Host

			pr.Out.Header.Set("X-Slot", "1")
			pr.Out.Header.Set("X-Forwarded-Proto", "https")
			pr.Out.Header.Set("X-Forwarded-Port", "443")
			pr.SetXForwarded()

			dump, _ := httputil.DumpRequestOut(pr.Out, true)
			fmt.Println(string(dump))

			fmt.Fprintf(os.Stdout,
				"Encaminhando requisição para o Traefik -> Host: %s | Headers: %s\n",
				pr.Out.Host, pr.Out.Header,
			)
		},
	}
	return &BlueGreen{
		next:  next,
		proxy: proxy,
		name:  name,
	}, nil
}
