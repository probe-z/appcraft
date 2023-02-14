package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// RequestContext
type RequestContext struct {
	Method  string        // http方法
	Timeout time.Duration // 超时时间
	Url     string        // url
	Request struct {
		Headers map[string]string // 请求头
		Body    []byte            // 请求包体
	}
	Err      error          // 错误信息
	TimeCost time.Duration  // 耗时
	Response *http.Response // 响应
}

// RequestContext Do
func (r *RequestContext) Do(ctx context.Context) {
	if r.Method == "" || r.Url == "" {
		r.Err = fmt.Errorf("bad request")
		return
	}
	beginTime := time.Now()
	subctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(subctx, r.Method, r.Url, bytes.NewBuffer(r.Request.Body))
	if err != nil {
		r.Err = err
		log.Printf("%v %v NewRequest err:%v\n", r.Method, r.Url, r.Err)
		return
	}
	for key, value := range r.Request.Headers {
		req.Header.Add(key, value)
	}
	req.Host = r.Request.Headers["host"]
	resp, err := gHTTPClient.Do(req)
	if err != nil || resp == nil {
		if err != nil {
			r.Err = err
		} else if resp == nil {
			r.Err = fmt.Errorf("NilResponse")
		}
		log.Printf("%v %v http.Do err:%v resp:%+v\n", r.Method, r.Url, r.Err, resp)
		return
	}
	r.TimeCost = time.Now().Sub(beginTime)
	r.Response = resp
	log.Printf("%v %v %v %v", r.Method, r.Url, r.TimeCost, resp.Header)
	return
}

// http配置
type HttpTransportConf struct {
	DialTimeout time.Duration `default:"30s"`
	KeepAlive   time.Duration `default:"30s"`

	// maximum amount of time waiting to wait for a TLS handshake
	TLSHandshakeTimeout time.Duration `default:"10s"`

	// maximum number of idle (keep-alive)
	// connections across all hosts.
	MaxIdleConns int `default:"100"`

	// maximum idle (keep-alive) connections to keep per-host.
	MaxIdleConnsPerHost int `default:"2"`

	// total number of connections per host, including connections in the dialing,
	// active, and idle states. On limit violation, dials will block.
	MaxConnsPerHost int `default:"0"`

	// maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing.
	IdleConnTimeout time.Duration `default:"90s"`

	// time to wait for a server's first response r.ReqHeaders after fully
	// writing the request r.ReqHeaders if the request has an
	// "Expect: 100-continue" header. Zero means no timeout and
	// causes the body to be sent immediately, without
	// waiting for the server to approve.
	// This time does not include the time to send the request header.
	ExpectContinueTimeout time.Duration `default:"1s"`

	// write buffer
	WriteBufferSize int `default:"4096"`

	// read buffer
	ReadBufferSize int `default:"4096"`

	// ForceAttemptHTTP2 controls whether HTTP/2 is enabled when a non-zero
	// Dial, DialTLS, or DialContext func or TLSClientConfig is provided.
	// By default, use of any those fields conservatively disables HTTP/2.
	// To use a custom dialer or TLS config and still attempt HTTP/2
	// upgrades, set this to true.
	ForceAttemptHTTP2 bool `default:"true"`
}

// The Client's Transport typically has internal state (cached TCP connections),
// so Clients should be reused instead of created as needed.
// Clients are safe for concurrent use by multiple goroutines.
var gHTTPClient *http.Client

// Configure
func Configure(conf *HttpTransportConf) {
	if gHTTPClient != nil {
		return
	}
	gHTTPClient = &http.Client{}
	gHTTPClient.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   conf.DialTimeout,
			KeepAlive: conf.KeepAlive,
			DualStack: true,
		}).DialContext,
		TLSHandshakeTimeout:   time.Duration(conf.TLSHandshakeTimeout),
		MaxIdleConns:          conf.MaxIdleConns,
		MaxIdleConnsPerHost:   conf.MaxIdleConnsPerHost,
		MaxConnsPerHost:       conf.MaxConnsPerHost,
		IdleConnTimeout:       time.Duration(conf.IdleConnTimeout),
		ExpectContinueTimeout: time.Duration(conf.ExpectContinueTimeout),
		WriteBufferSize:       conf.WriteBufferSize,
		ReadBufferSize:        conf.ReadBufferSize,
		ForceAttemptHTTP2:     conf.ForceAttemptHTTP2,
	}
}
