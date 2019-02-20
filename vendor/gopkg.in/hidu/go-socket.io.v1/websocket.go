package socketio
import "github.com/equalll/mydebug"

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"io"
	"time"
)

func init() {mydebug.INFO()
	DefaultTransports.RegisterTransport("websocket")
}

type webSocket struct {
	session *Session
	timeout time.Duration
	conn    *websocket.Conn
}

func newWebSocket(session *Session) *webSocket {mydebug.INFO()
	ret := &webSocket{
		session: session,
		timeout: session.heartbeatTimeout / 10,
	}
	session.transport = ret
	return ret
}

func (ws *webSocket) Send(data []byte) error {mydebug.INFO()
	ws.conn.SetWriteDeadline(time.Now().Add(ws.timeout))
	return websocket.Message.Send(ws.conn, string(data))
}

func (ws *webSocket) Read() (io.Reader, error) {mydebug.INFO()
	var ret string
	ws.conn.SetReadDeadline(time.Now().Add(ws.timeout))
	err := websocket.Message.Receive(ws.conn, &ret)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewBufferString(ret)
	return reader, nil
}

func (ws *webSocket) webSocketHandler(conn *websocket.Conn) {mydebug.INFO()
	ws.conn = conn
	ws.session.loop()
}
