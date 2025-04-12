package server

import (
	"fmt"
	"net"
)

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
