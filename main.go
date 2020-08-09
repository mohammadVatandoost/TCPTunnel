package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

func main() {

	// connection1 := &Endpoint{
	// 	Host: "localhost",
	//     Port: 4567,
	// }

	connection1 := &Endpoint{
		Host: "192.168.2.9",
		Port: 22,
	}

	connection2 := &Endpoint{
		Host: "localhost",
		Port: 4568,
	}

	// listener, err := net.Listen("tcp", connection1.String())
	// if err != nil {
	// 	fmt.Println("connection1 error :", err)
	// 	os.Exit(1)
	// }

	conn1, err := net.Dial("tcp", connection1.String())
	if err != nil {
		fmt.Printf("connection error: %s", err)
		os.Exit(1)
	}

	listener2, err := net.Listen("tcp", connection2.String())
	if err != nil {
		fmt.Println("connection2 error :", err)
		os.Exit(1)
	}

	// defer listener.Close()
	defer conn1.Close()

	defer listener2.Close()

	for {
		// conn1, err := listener.Accept()
		// if err != nil {
		// 	fmt.Println("error :", err)
		//     os.Exit(1)
		// }
		// fmt.Println("connection 1 accept ")
		conn2, err := listener2.Accept()
		if err != nil {
			fmt.Println("error :", err)
			os.Exit(1)
		}
		fmt.Println("connection 2 accept ")
		go copyConn(conn1, conn2)
		go copyConn(conn2, conn1)

	}
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
