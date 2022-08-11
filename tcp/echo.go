package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/kun98-liu/MyGodis/lib/logger"
	"github.com/kun98-liu/MyGodis/lib/sync/atomic"
	"github.com/kun98-liu/MyGodis/lib/sync/wait"
)

//implement the interface Handler
type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		_ = conn.Close()
		return
	}

	client := &EchoClient{
		Conn: conn,
	}

	h.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("Read EOF")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}

		client.Waiting.Add(1)
		b := []byte(msg)
		conn.Write(b)
		client.Waiting.Done()
	}

}

func (h *EchoHandler) Close() error {
	logger.Info("EchoHandler is shutting down...")
	h.closing.Set(true)
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(*EchoClient)
		client.Close()
		return true
	})
	return nil
}

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (c *EchoClient) Close() {
	c.Waiting.WaitWithTimeout(10 * time.Second)
	logger.Info("EchoClient is shutting down...")
	c.Conn.Close()
}
