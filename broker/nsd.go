package broker

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/nats-io/go-nats"
)

const (
	// SubRegister nats topic of
	SubRegister = "nsd.register"
)

// ServerCfg configuration of nats service discover server.
type ServerCfg struct {
	IP   string
	Port int
	Nats string
}

// Server nsd server.
type Server struct {
	Cfg  ServerCfg
	Conn *nats.Conn
}

// Start nsd server
func (s Server) Start() (err error) {
	s.Conn, err = nats.Connect(
		s.Cfg.Nats,
		nats.Name("nsd client"),
		nats.Timeout(4*time.Second),
		// nats.ReconnectHandler(),
	)
	if err != nil {
		return
	}
	// s.Conn.Subscribe()
	return
}

// LocalAddress list local addresses.
func LocalAddress() string {
	list, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for i, iface := range list {
		fmt.Printf("%d name=%s %v\n", i, iface.Name, iface)
		if strings.HasPrefix(iface.Name, "lo") {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			panic(err)
		}
		for j, addr := range addrs {
			fmt.Printf("- %d %v\n", j, addr)
		}
	}
	return ""
}
