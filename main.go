package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

func main() {
	var wg sync.WaitGroup
	// connection1 := &Endpoint{
	// 	Host: "localhost",
	// 	Port: 4570,
	// }

	connection1 := &Endpoint{
		Host: "0.0.0.0",
		Port: 4570,
	}

	// connection1 := &Endpoint{
	// 	Host: "192.168.2.9",
	// 	Port: 22,
	// }

	connection2 := &Endpoint{
		Host: "0.0.0.0",
		Port: 4567,
	}

	listener1, err := net.Listen("tcp", connection1.String())
	if err != nil {
		fmt.Println("connection1 error :", err)
		os.Exit(1)
	}

	listener2, err := net.Listen("tcp", connection2.String())
	if err != nil {
		fmt.Println("connection2 error :", err)
		os.Exit(1)
	}

	// defer listener1.Close()
	// defer listener2.Close()

	// channelL1 := make(chan net.Conn)
	// channelL2 := make(chan net.Conn)
	// wg.Add(1)
	// go conListner(listener1, channelL1, channelL2, &wg)
	// wg.Add(1)
	// go conListner(listener2, channelL2, channelL2, &wg)

	for {

		conn1, err := listener1.Accept()
		if err != nil {
			fmt.Println("error :", err)
			os.Exit(1)
		}
		fmt.Println("device connection  accepted ")
		go connectionListener(listener2, conn1)
		// go copyConn(conn1, conn2)
		// go copyConn(conn2, conn1)

	}

	// wg.Wait()
	fmt.Println("Main: Completed")
}

func connectionListener(listener net.Listener, conn1 net.Conn) {
	for {
		conn2, err := listener.Accept()
		if err != nil {
			fmt.Println("error :", err)
			os.Exit(1)
		}
		fmt.Println("connection 2 accept ")
		go copyConn(conn1, conn2)
		go copyConn(conn2, conn1)

	}
}

func conListner(listener net.Listener, channelTx chan net.Conn, channelRx chan net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		conn1, err := listener.Accept()
		if err != nil {
			fmt.Println("error :", err)
			os.Exit(1)
		}
		fmt.Println("connection  accept ")
		channelTx <- conn1
		conn2 := <-channelRx
		go copyConn(conn1, conn2)
	}

	defer listener.Close()
}

func copyConn(writer, reader net.Conn) {
	defer writer.Close()
	defer reader.Close()

	_, err := io.Copy(writer, reader)
	if err != nil {
		fmt.Printf("io.Copy error: %s", err)
	}

	fmt.Println("copyConn")
}

// func tunnel(conn1 net.Conn, conn2 net.Conn) {

// }
