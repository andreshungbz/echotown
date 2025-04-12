// Package server contains functions to start and manage the echo server.
package server

import (
	"bufio"
	"fmt"
	"net"

	"github.com/andreshungbz/echotown/internal/logger"
)

// Start launches an infinite loop that creates a goroutine for every connecting client.
func Start(port int) {
	// TCP connection listener on the local machine
	listener, err := createTCPListener(port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	// logger to keep track of connections and disconnections
	serverLogger, close, err := logger.NewServer()
	if err != nil {
		panic(err)
	}
	defer close()

	// print and log server start message
	serverLogger.Printf("[INFO] Echo Town server started at [%s:%d]\n", getLocalIP(), port)

	// monitor for termination signal and print server stop message
	monitorTermSig(serverLogger)

	// server infinite loop
	for {
		// wait to accept incoming client connections
		conn, err := listener.Accept()
		if err != nil {
			serverLogger.Printf("[ERROR] Connection acceptance failed: %v\n", err)
			continue
		}

		// launch a goroutine to handle the individual client
		serverLogger.Printf("[INFO] [%v] connected to the server.\n", conn.RemoteAddr())
		go func() {
			handleConn(conn)
			serverLogger.Printf("[INFO] [%v] disconnected from the server.\n", conn.RemoteAddr())
		}()

		// loop back to wait and accept another client connection until Ctrl + C is pressed on server
	}
}

// handleConn processes an individual connection to a client.
func handleConn(conn net.Conn) {
	defer conn.Close()

	clientAddress := conn.RemoteAddr()

	// send a welcome message to the client
	_, err := conn.Write([]byte(fmt.Sprintf("Welcome to Echo Town! (CTRL + C to disconnect)\nYou are connected as [%v]\n", clientAddress)))
	if err != nil {
		fmt.Println("Error sending welcome message to client:", err)
		return
	}

	reader := bufio.NewReader(conn)

	// client connection infinite loop
	for {
		// indicate to client a prompt for input
		_, err = conn.Write([]byte(fmt.Sprintf("\n[%v]: ", clientAddress)))
		if err != nil {
			fmt.Println("Error sending prompt to client:", err)
			return
		}

		response, err := createResponse(reader)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}

		// write response to the client
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to client:", err)
		}
	}
}
