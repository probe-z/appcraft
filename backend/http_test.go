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

type HiRsp struct {
	Name string
}

func Test(t *testing.T) {
	httpClientConf := &httpclient.HttpTransportConf{}
	defaults.Set(httpClientConf)
	httpclient.Configure(httpClientConf)

	r := &httpclient.RequestContext{
		Method:  "GET",
		Timeout: time.Duration(2) * time.Second,
		Url:     "http://www.baidu.com",
	}
	r.Do(context.Background())
	svrConf := httpserver.ServerConfig{
		Addr:           "127.0.0.1:8080",
		RestartTimeout: 3,
	}
	go func() {
		time.Sleep(time.Duration(5) * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR2)
	}()
	srv := httpserver.NewServer(&svrConf)
	srv.HandleFunc("/hi", func(r httpserver.RequestContext) {
		log.Printf("hi")
		r.SetJsonResponse(&httpserver.JsonResponse{
			Data: &HiRsp{
				Name: "dddd",
			},
		})
	})
	srv.Run()
}
