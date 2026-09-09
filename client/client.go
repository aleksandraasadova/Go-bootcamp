package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const addr string = ":8080"

func main() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	readBytes, err := io.ReadAll(conn)
	if err != nil && err != io.EOF {
		fmt.Printf("Failed to read response: %v\n", err)
		os.Exit(1)
	}
	if string(readBytes) == "OK\n" {
		fmt.Println("Success: received correct response 'OK\\n'")
	} else {
		fmt.Printf("Error: expected 'OK\\n', but got '%s'\n", string(readBytes))
	}
}
