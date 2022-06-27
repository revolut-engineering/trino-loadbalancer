package main

import (
	"log"
	"net"
	"net/url"
	"time"
)

// pingBackend checks if the backend is alive.
func isAlive(url *url.URL) bool {
	conn, err := net.DialTimeout("tcp", url.Host, time.Minute*1)

	if err != nil {
		log.Printf("Unreachable to %v, error: %s", url.Host, err.Error())
		return false
	}

	defer conn.Close()
	return true
}

// healthCheck is a function for health check
func healthCheck(cfg *Config) {
	t := time.NewTicker(time.Second * time.Duration(cfg.HealthCheck.Period))
	for {
		select {
		case <-t.C:
			for i := 0; i < len(cfg.Backends); i++ {
				backend := &cfg.Backends[i]

				pingURL, err := url.Parse(backend.URL)
				if err != nil {
					log.Fatal(err.Error())
				}

				isAlive := isAlive(pingURL)
				backend.SetDead(!isAlive)
				msg := "ok"
				if !isAlive {
					msg = "dead"
				}

				log.Printf("%v checked %v by health check", backend.URL, msg)
			}
		}
	}

}
