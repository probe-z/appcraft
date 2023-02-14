package main

import (
	"context"
	"github.com/creasty/defaults"
	"httpclient"
	"httpserver"
	"log"
	"syscall"
	"testing"
	"time"
)

func Test(t *testing.T) {
	// test server
	go func() {
		svrConf := httpserver.ServerConfig{
			Addr:           "127.0.0.1:8080",
			RestartTimeout: 3,
		}
		srv := httpserver.NewServer(&svrConf)
		srv.HandleFunc("/hi", func(r httpserver.RequestContext) {
			log.Printf("hi")
			r.SetJsonResponse(&httpserver.JsonResponse{
				Data: &struct{ Name string }{
					Name: "dddd",
				},
			})
		})
		srv.Run()
	}()

	// test client
	time.Sleep(time.Duration(1) * time.Second)
	httpClientConf := &httpclient.HttpTransportConf{}
	defaults.Set(httpClientConf)
	httpclient.Configure(httpClientConf)
	(&httpclient.RequestContext{
		Method:  "GET",
		Timeout: time.Duration(2) * time.Second,
		Url:     "http://127.0.0.1:8080/hi",
	}).Do(context.Background())
	syscall.Kill(syscall.Getpid(), syscall.SIGUSR2)
}
