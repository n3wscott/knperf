package main

import (
	"flag"
	"github.com/kelseyhightower/envconfig"
	"io"
	"log"
	"net/http"
	"os"
)

type Proxy struct {
	Target string
}

func NewProxy(target string) *Proxy {
	return &Proxy{Target: target}
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, p.Target, r.Body) // TODO: maybe pull the path info.
	for k, vs := range r.Header {
		for _, v := range vs {
			req.Header.Set(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = r.Body.Close()

	for k, vs := range resp.Header {
		for _, v := range vs {
			rw.Header().Set(k, v)
		}
	}
	rw.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(rw, resp.Body)
	_ = resp.Body.Close()
}

type envConfig struct {
	// Target URL where to proxy
	Target string `envconfig:"TARGET" required:"true"`
}

func main() {
	flag.Parse()

	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}
	log.Fatal(http.ListenAndServe(":8080", NewProxy(env.Target)))
}
