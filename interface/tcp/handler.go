package tcp

import (
	"context"
	"net"
)

//HandlerFunc is the application handler function
type HandlerFunc func(ctx context.Context, conn net.Conn)

//A handler is an application server built on tcp
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}
