package server

import (
	"fmt"
	"net"

	"github.com/andreshungbz/echotown/internal/logger"
)

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

	fmt.Println("Server listening on :4000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}

		serverLogger.Printf("[%v] connected to the server.\n", conn.RemoteAddr())

		go func() {
			handleConn(conn)
			serverLogger.Printf("[%v] disconnected from the server.\n", conn.RemoteAddr())
		}()
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

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
