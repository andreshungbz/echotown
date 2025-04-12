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

// getLocalAddr returns [net.Addr] of the local machine with a passed in port.
func getLocalAddr(port int) net.Addr {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return &net.TCPAddr{IP: net.ParseIP("localhost"), Port: port}
	}

	for _, addr := range addrs {
		// check that the address is not the loopback address
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// return only IPv4 addresses
			if ipNet.IP.To4() != nil {
				return &net.TCPAddr{IP: ipNet.IP, Port: port}
			}
		}
	}

	return &net.TCPAddr{IP: net.ParseIP("localhost"), Port: port}
}
