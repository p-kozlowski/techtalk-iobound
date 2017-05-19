// source: https://coderwall.com/p/wohavg/creating-a-simple-tcp-server-in-go

package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
	"sync/atomic"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
	DEFAULT_DELAY = time.Duration(5) * time.Millisecond
)

var activeCount int32 = 0

func main() {
	listenHost := flag.String("listen-host", CONN_HOST, "IP address for listening")
	port := flag.String("port", CONN_PORT, "port for listening")
	delay := flag.Duration("delay", DEFAULT_DELAY, "response delay")

	flag.Parse()

	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, *listenHost+":"+*port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	go printActive()

	fmt.Println("Listening on " + *listenHost + ":" + *port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, *delay)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, delay time.Duration) {
	defer conn.Close()
	atomic.AddInt32(&activeCount, 1)
	defer atomic.AddInt32(&activeCount, -1)

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err.Error(), message)
		return
	}
	time.Sleep(delay)

	// Send a response back to person contacting us.
	fmt.Fprint(conn, message)
}

func printActive() {
	for _ = range time.Tick(time.Second) {
		fmt.Println("Active connections: ", atomic.LoadInt32(&activeCount))
	}
}