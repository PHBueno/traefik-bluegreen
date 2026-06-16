package traefik_bluegreen

import (
	"context"
	"fmt"
	"net/http"
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
	next http.Handler
	name string
}

func (bg *BlueGreen) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	req.Header.Set("X-Slot", "1")
	fmt.Fprintln(os.Stdout, "TESTE => ", req)
	fmt.Fprintln(os.Stdout, "Chamando o ServeHTTP")

	bg.next.ServeHTTP(rw, req)
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.RedisAddress == "" {
		return nil, fmt.Errorf("Redis Address is not set!")
	}
	return &BlueGreen{
		next: next,
		name: name,
	}, nil
}
