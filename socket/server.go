package socket

import (
	"github.com/tidwall/evio"
)

// Server handle connections
type Server struct {
	Server evio.Server
}
