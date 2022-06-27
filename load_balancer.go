package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
)

type loadBalancer struct {
	mu           sync.Mutex
	idx          int
	cfg          *Config
	job_registry map[string]int
}

func (l *loadBalancer) lbHandler(w http.ResponseWriter, r *http.Request) {
	var targetURL *url.URL
	var currentBackend *Backend
	var current_id int

	if r.Method == "POST" {
		l.mu.Lock()
		for {
			current_id = l.idx % len(l.cfg.Backends)
			l.idx++

			currentBackend = &l.cfg.Backends[current_id]
			if currentBackend.isAlive() {
				break
			} else {
				log.Printf("Backend %s is not alive", currentBackend.URL)
				currentBackend.SetDead(true)
			}
		}

		targetURL = l.get_target_url(current_id)
		l.mu.Unlock()
	} else {
		status := status_from_url(r.URL)
		job_id := l.job_registry[status.job_id]

		targetURL = l.get_target_url(job_id)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)

	if r.Method == "POST" {
		reverseProxy.ModifyResponse = func(resp *http.Response) (err error) {
			b, _ := ioutil.ReadAll(resp.Body)
			body := parse_body(b)
			l.job_registry[body.ID] = current_id

			resp.Body = ioutil.NopCloser(bytes.NewReader(b))
			return nil
		}
	}

	reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		log.Printf("%v is dead.", targetURL)
		currentBackend.SetDead(true)

		l.lbHandler(w, r)
	}

	reverseProxy.ServeHTTP(w, r)
}

func (l *loadBalancer) get_target_url(id int) *url.URL {
	targetURL, err := url.Parse(l.cfg.Backends[id].URL)
	if err != nil {
		log.Fatal(err.Error())
	}

	return targetURL
}

type queryStatus struct {
	status string
	job_id string
}

func status_from_url(url *url.URL) queryStatus {
	split := strings.Split(url.Path, "/")

	// Path starts with / therefore increasing the array index by 1
	return queryStatus{
		status: split[3],
		job_id: split[4],
	}
}
