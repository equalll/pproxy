package socketio
import "github.com/equalll/mydebug"

import (
	"io"
	"sync"
)

type Transport interface {
	Send([]byte) error
	Read() (io.Reader, error)
}

var (
	DefaultTransports = NewTransportManager()
)

type TransportManager struct {
	mutex      sync.RWMutex
	transports map[string]bool
}

func NewTransportManager() *TransportManager {mydebug.INFO()
	return &TransportManager{
		transports: make(map[string]bool),
	}
}

func (tm *TransportManager) RegisterTransport(name string) {mydebug.INFO()
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tm.transports[name] = true
}

func (tm *TransportManager) GetTransportNames() (names []string) {mydebug.INFO()
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	names = make([]string, 0, len(tm.transports))
	for k, _ := range tm.transports {
		names = append(names, k)
	}
	return
}
