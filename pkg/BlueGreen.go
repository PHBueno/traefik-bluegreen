package pkg

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
)

type BlueGreen struct {
	next  http.Handler
	proxy *httputil.ReverseProxy
	name  string
}

func (bg *BlueGreen) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	slog.Info("requisição chegou", "req", req)
	// Evita loop
	if req.Header.Get("X-Slot") != "" {
		http.NotFound(rw, req)
		return
	}

	bg.proxy.ServeHTTP(rw, req)

}

func New(next http.Handler, proxy *httputil.ReverseProxy, name string) *BlueGreen {
	return &BlueGreen{
		next:  next,
		proxy: proxy,
		name:  name,
	}
}
