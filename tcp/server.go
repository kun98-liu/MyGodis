package tcp

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kun98-liu/MyGodis/interface/tcp"
	"github.com/kun98-liu/MyGodis/lib/logger"
)

type Config struct {
	Address    string        `yaml:"address"`
	MaxConnect uint32        `yaml:"max-connect"`
	Timeout    time.Duration `yaml:"timeout"`
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	signalChan := make(chan os.Signal)

	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		signal := <-signalChan
		switch signal {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", cfg.Address)

	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("bind: %s, start listening", cfg.Address))
	ListenAndServe(listener, handler, closeChan)
	return nil

}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan chan struct{}) {

	go func() {
		<-closeChan
		logger.Info("shutting down...")
		_ = listener.Close()
		_ = handler.Close()
	}()

	defer func() {
		_ = listener.Close()
		_ = handler.Close()
		logger.Info("shutting down anyway...")
	}()

	ctx := context.Background()

	var waitDone sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("Accept link")
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
				logger.Info("Handler finished its task")
			}()
			logger.Info("Handler start to handle task")
			handler.Handle(ctx, conn)
		}()
	}

	waitDone.Wait()

}
