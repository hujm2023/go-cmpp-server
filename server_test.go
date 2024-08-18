package cmppserver

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewCMPPServer()
	errChan := make(chan error)
	go func() {
		errChan <- s.Listen("tcp", "0.0.0.0:8899")
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	select {
	case s := <-signals:
		fmt.Printf("receive signal: %d", s)
		return
	case err := <-errChan:
		fmt.Printf("run server error: %v", err)
	}
}
