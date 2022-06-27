package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HealthCheck HealthCheck `yaml:"health_check"`
	Proxy       Proxy       `yaml:"proxy"`
	Backends    []Backend   `yaml:"backends"`
}

type Proxy struct {
	Port string `yaml:"port"`
}

type Backend struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url"`
	IsDead bool
	mu     sync.RWMutex
}

type HealthCheck struct {
	Status bool `yaml:"status"`
	Period int  `yaml:"period"`
}

func (backend *Backend) SetDead(b bool) {
	backend.mu.Lock()
	backend.IsDead = b
	backend.mu.Unlock()
}

func (backend *Backend) isAlive() bool {
	if backend.IsDead {
		return false
	}

	pingURL, err := url.Parse(backend.URL)
	if err != nil {
		log.Fatal(err.Error())
	}

	conn, err := net.DialTimeout("tcp", pingURL.Host, time.Minute*1)

	if err != nil {
		log.Printf("Unreachable to %v, error: %s", pingURL.Host, err.Error())
		backend.SetDead(true)

		return false
	}

	defer conn.Close()
	return true
}

func read_config(name string) Config {
	var cfg Config

	data, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatal(err.Error())
	}

	err2 := yaml.Unmarshal(data, &cfg)
	if err2 != nil {
		log.Fatal(err2)
	}

	return cfg
}
