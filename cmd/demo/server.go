package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	cmppserver "github.com/hujm2023/go-cmpp-server"
)

func main() {
	s := cmppserver.NewCMPPServer()
	errChan := make(chan error)
	go func() {
		errChan <- s.Listen("tcp", "[::]:8899")
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	select {
	case s := <-signals:
		fmt.Printf("receive signal: %d\n", s)
		return
	case err := <-errChan:
		fmt.Printf("run server error: %v\n", err)
	}
}
