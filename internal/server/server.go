// Package server contains functions to start and manage the echo server.
package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/andreshungbz/echotown/internal/server/internal/command"
	"github.com/andreshungbz/echotown/internal/server/internal/logger"
	"github.com/andreshungbz/echotown/internal/server/internal/personality"

	"github.com/fatih/color"
)

// Start infinitely listens for connection and creates a goroutine for every connecting client.
// It is stopped when an interrupt or termination signal is received.
func Start(port int) { // launch a goroutine to handle the individual client
	// create TCP listener on specified port
	listener, err := createTCPListener(port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	// logger for connections/disconnections
	serverLogger, close, err := logger.NewServer()
	if err != nil {
		panic(err)
	}
	defer close()

	// monitor for termination signal and print server stop message
	monitorTermSig(serverLogger)

	serverLogger.Printf("[INFO] Echo Town server started at [%v]\n", getLocalAddr(port))
	for { // server infinite loop
		conn, err := listener.Accept() // wait to accept incoming client connections
		if err != nil {
			serverLogger.Printf("[ERROR] Connection acceptance failed: %v\n", err)
			continue
		}

		clientConnectLog := fmt.Sprintf("[INFO] [%v] connected to the server.\n", conn.RemoteAddr())
		clientDisconnectLog := fmt.Sprintf("[INFO] [%v] disconnected from the server.\n", conn.RemoteAddr())

		serverLogger.Print(clientConnectLog)
		// launch a goroutine to handle the individual client
		go func() {
			// logger for client messages/server responses
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
	}
}

// handleConn processes an individual connection to a client.
func handleConn(conn net.Conn, serverLogger, clientLogger *log.Logger) {
	defer conn.Close()

	clientAddress := conn.RemoteAddr()
	clientPrompt := color.YellowString("\n[%v]: ", clientAddress)
	serverPrepend := color.CyanString("[Echo Town]: ")
	welcomeMessage := fmt.Sprintf("Welcome to Echo Town!\nEnter /quit or \"bye\" to exit\nEnter /help for more commands\nYou are connected as [%v]\n", clientAddress)
	goodbyeMessage := "\nCome back to Echo Town soon!\n"

	_, err := conn.Write([]byte(welcomeMessage)) // send a welcome message to the client
	if err != nil {
		serverLogger.Print(createError("Welcome Message", clientAddress, err))
		return
	}

	reader := bufio.NewReader(conn)

	for { // client connection infinite loop
		_, err = conn.Write([]byte(clientPrompt)) // indicate to client a prompt for input
		if err != nil {
			serverLogger.Print(createError("Client Prompt", clientAddress, err))
			return
		}

		conn.SetReadDeadline(time.Now().Add(time.Second * 30)) // set or reset connection timeout

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

		// PERSONALITY & COMMAND PROTOCOL PARSING

		var close bool // determines whether to close the connection based on command protocol or personality

		response, close = command.Parse(response)

		if close {
			conn.Write([]byte(goodbyeMessage))
			return
		}

		response, close = personality.Parse(response)

		// SERVER TO CLIENT WRITING

		clientLogger.Printf("[RESPONSE] %s", response)            // log server response to client log
		response = fmt.Sprintf("%s%s\n", serverPrepend, response) // prepend server response

		_, err = conn.Write([]byte(response)) // write response to the client
		if err != nil {
			serverLogger.Print(createError("Server Write", clientAddress, err))
		}

		if close {
			conn.Write([]byte(goodbyeMessage))
			return
		}
	}
}
