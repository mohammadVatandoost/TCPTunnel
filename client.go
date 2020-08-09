// socket client for golang
// https://golangr.com
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}
	portNum := os.Args[1]
	// connect to server
	conn, err := net.Dial("tcp", "127.0.0.1:"+portNum)
	if err != nil {
		fmt.Printf("connection error: %s", err)
		os.Exit(1)
	}

	go receiveFromClient(conn)
	sendToServer(conn)
}

func sendToServer(conn net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		// send to server
		fmt.Fprintf(conn, text+"\n")
	}
}

func receiveFromClient(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		checkError(err)
		fmt.Println("Message from client: " + message)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
