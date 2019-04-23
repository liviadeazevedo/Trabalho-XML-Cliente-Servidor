package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "4446"
	CONN_TYPE = "tcp4"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	for {
		// Make a buffer to hold incoming data.
		buf := make([]byte, 8096)
		// Read the incoming connection into the buffer.
		reqLen, err := bufio.NewReader(conn).Read(buf) //conn.Read(buf)

		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}

		fmt.Println("Número de bytes lidos do cliente: ",reqLen)
		fmt.Println("Array de bytes lido convertido para string:\n\n ",string(buf))

		// Send a response back to person contacting us.
		conn.Write([]byte("Message received."))
	}
	// Close the connection when you're done with it.
	conn.Close()
}