package httpserver

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server
type Server struct {
	*http.Server
	handler *mux.Router
	conf    *ServerConfig
}

// NewServer
func NewServer(conf *ServerConfig) *Server {
	return &Server{
		conf: conf,
	}
}

// RequestContext
type RequestContext struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

// HandleFunc
func (s *Server) HandleFunc(path string, f func(RequestContext)) {
	if s.handler == nil {
		s.handler = mux.NewRouter()
	}
	s.handler.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		f(RequestContext{
			Request:        r,
			ResponseWriter: w,
		})
	})
}

// Run
func (s *Server) Run() {
	if s.Server == nil {
		s.Server = &http.Server{
			Addr:           s.conf.Addr,
			ReadTimeout:    time.Second * time.Duration(s.conf.ReadTimeout),
			WriteTimeout:   time.Second * time.Duration(s.conf.WriteTimeout),
			IdleTimeout:    time.Second * time.Duration(s.conf.IdleTimeout),
			MaxHeaderBytes: s.conf.MaxHeaderBytes,
			Handler:        s.handler,
		}
	}
	go func() {
		if err := s.Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR2)
	<-ch

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.conf.RestartTimeout)*time.Second)
	defer cancel()
	s.Server.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
}

// ServerConfig
type ServerConfig struct {
	Addr           string
	ReadTimeout    int
	WriteTimeout   int
	IdleTimeout    int
	MaxHeaderBytes int
	RestartTimeout int
}
