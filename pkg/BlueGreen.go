package pkg

import (
	"net/http"
	"net/http/httputil"
)

type BlueGreen struct {
	next  http.Handler
	proxy *httputil.ReverseProxy
	name  string
}

func (bg *BlueGreen) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Evita loop
	if req.Header.Get("X-Slot") != "" {
		bg.next.ServeHTTP(rw, req)
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
