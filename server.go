package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const port string= ":8080"

type tcpServer struct {
	addr     string
	listener net.Listener
}

func newTCPServer(addr string) *tcpServer {
	return &tcpServer{
		addr: addr,
	}
}

func (s *tcpServer) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to start TCP server on port %s: %v", s.addr, err)
	}
	s.listener = listener

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if err == net.ErrClosed {
				return nil
			}
			fmt.Printf("failed to connect the client: %v\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *tcpServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *tcpServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	_, err := conn.Write([]byte("OK\n"))
	if err != nil {
		fmt.Printf("Failed to send the connection message: %v\n", err)
	}
}

func main() {
	server := newTCPServer(port)
	go func() {
		if err := server.Start(); err != nil {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
	<-sigint
	server.Stop()
}
