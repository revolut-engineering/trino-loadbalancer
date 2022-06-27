package main

import (
	"log"
	"net/http"
)

// Serve serves a loadbalancer.
func main() {
	cfg := read_config("./config.yaml")
	log.Println(cfg)

	if cfg.HealthCheck.Status {
		go healthCheck(&cfg)
	}

	load_balancer := loadBalancer{
		cfg:          &cfg,
		job_registry: make(map[string]int),
	}

	s := http.Server{
		Addr:    ":" + cfg.Proxy.Port,
		Handler: http.HandlerFunc(load_balancer.lbHandler),
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
