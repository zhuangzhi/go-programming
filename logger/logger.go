package logger

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

// LogPipline pipline and filter log to other services.
type LogPiplineCfg struct {
	File      string
	Collector string
}
type Filter struct {
}

// Collector c
type Agent struct {
	ServerURL string
}

// Start is
func (c Agent) Start() error {
	r, w := io.Pipe()
	out := os.Stdout
	os.Stdout = w
	go func() {
		for {
			br := bufio.NewReader(r)
			br.ReadLine()
		}
	}()
}

type Server struct {
	Port int
}

// Start is
func (s *Server) Start() {
	s.router = httpRouter.New()
	s.router.POST("/", s.handler)

	server := &http.Server{
		Handler: router,
	}
	listener, err := net.Listen("tcp", ":10000")
	if err != nil {
		return err
	}
	log.Println("HTTP server is listening..")
	return server.ServeTLS(
		listener,
		"./certs/key.crt",
		"./certs/key.key",
	)
}

func (s *Server) handler(
	w http.ResponseWriter,
	req *http.Request,
	_ httpRouter.Params,
) {
	// Cert
}
