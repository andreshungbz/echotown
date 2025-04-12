// Package server contains functions to start the echo server.
package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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

// NON-EXPORTED FUNCTIONS

// handleConn processes an individual connection to a client.
func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	// client connection infinite loop
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}

		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("Error writing to client:", err)
		}
	}
}

// createTCPListener returns a [net.Listener] with the network set to "tcp"
// and the address set to the local machine along with the provided port argument.
func createTCPListener(port int) (net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	return listener, nil
}

// getLocalIP returns the local IPv4 address string or "localhost"
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}

	for _, addr := range addrs {
		// check that the address is not the loopback address
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// return only IPv4 addresses
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return "localhost"
}

// monitorTermSig creates a channel to monitor OS termination signals
// and launches a single goroutine to print and log a server end message when
// the signal is received.
func monitorTermSig(logger *log.Logger) {
	signalChan := make(chan os.Signal, 1)

	// register channel to receive interrupt and termination OS signals
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// when a proper signal is received, log message and exit program
	go func() {
		<-signalChan
		fmt.Println() // print a new line just for the console
		logger.Print("[INFO] Echo Town server stopped.\n\n")
		os.Exit(0)
	}()
}
