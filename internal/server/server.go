// Package server contains functions to start and manage the echo server.
package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

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
	serverLogger.Printf("[INFO] Echo Town server started at [%v]\n", getLocalAddr(port))

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

		clientConnectLog := fmt.Sprintf("[INFO] [%v] connected to the server.\n", conn.RemoteAddr())
		clientDisconnectLog := fmt.Sprintf("[INFO] [%v] disconnected from the server.\n", conn.RemoteAddr())

		// launch a goroutine to handle the individual client
		serverLogger.Print(clientConnectLog)
		go func() {
			// logger to record client messages and server responses to those messages
			clientLogger, close, err := logger.NewClient(conn.RemoteAddr())
			if err != nil {
				panic(err)
			}
			defer close()

			clientLogger.Print(clientConnectLog)
			handleConn(conn, serverLogger, clientLogger)
			clientLogger.Print(clientDisconnectLog + "\n")

			serverLogger.Print(clientDisconnectLog)
		}()

		// loop back to wait and accept another client connection until Ctrl + C is pressed on server
	}
}

// handleConn processes an individual connection to a client.
func handleConn(conn net.Conn, serverLogger, clientLogger *log.Logger) {
	defer conn.Close()

	clientAddress := conn.RemoteAddr()
	clientPrompt := fmt.Sprintf("\n[%v]: ", clientAddress)
	welcomeMessage := fmt.Sprintf("Welcome to Echo Town! (CTRL + C to disconnect)\nYou are connected as [%v]\n", clientAddress)

	// send a welcome message to the client
	_, err := conn.Write([]byte(welcomeMessage))
	if err != nil {
		serverLogger.Print(createError("Welcome Message", clientAddress, err))
		return
	}

	conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	reader := bufio.NewReader(conn)

	// client connection infinite loop
	for {
		// indicate to client a prompt for input
		_, err = conn.Write([]byte(clientPrompt))
		if err != nil {
			serverLogger.Print(createError("Client Prompt", clientAddress, err))
			return
		}

		// construct response from validated input
		response, err := createResponse(reader, clientLogger)
		if err != nil {
			errString := createError("Client Input", clientAddress, err)

			// print a nicer message for client when timeout occurs
			if strings.Contains(errString, "timeout") {
				conn.Write([]byte("\n[TIMEOUT] 30 seconds have passed. Disconnecting...\n"))
				serverLogger.Print(errString)
				return
			}

			conn.Write([]byte(errString))
			serverLogger.Print(errString)
			return
		}

		// write response to the client
		_, err = conn.Write([]byte(response))
		if err != nil {
			serverLogger.Print(createError("Server Write", clientAddress, err))
		}
	}
}
