package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

func main() {
	// var wg sync.WaitGroup
	// connection1 := &Endpoint{
	// 	Host: "localhost",
	// 	Port: 4570,
	// }
	if len(os.Args) < 4 {
		fmt.Println("Please Provide Port numbers")
		os.Exit(0)
	}
	var portNum1 int
	var portNum2 int
	var err error
	state := 0
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "-p1" && state == 0 {
			state = 1
		} else if state == 1 {
			portNum1, err = strconv.Atoi(os.Args[i])
			if err != nil {
				fmt.Println("Please Provice valid port number, err:", err)
				os.Exit(0)
			}
			state = 0
		} else if os.Args[i] == "-p2" && state == 0 {
			state = 2
		} else if state == 2 {
			portNum2, err = strconv.Atoi(os.Args[i])
			if err != nil {
				fmt.Println("Please Provice valid port number, err:", err)
				os.Exit(0)
			}
			state = 0
		}
	}

	connection1 := &Endpoint{
		Host: "0.0.0.0",
		Port: portNum1,
	}

	// connection1 := &Endpoint{
	// 	Host: "192.168.2.9",
	// 	Port: 22,
	// }

	connection2 := &Endpoint{
		Host: "0.0.0.0",
		Port: portNum2,
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
	fmt.Println("portNum1:", portNum1, ", portNum2: ", portNum2)
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
	cChannel1 := make(chan int)
	cChannel2 := make(chan int)
	sigChannel := make(chan int)
	// for {

	select {
	case msg1 := <-cChannel1:
		// conn2.Close()
		fmt.Println("conn2 closed :", msg1)
	case msg2 := <-cChannel2:
		fmt.Println("connection 1 closed :", msg2)
		return
	default:
		conn2, err := listener.Accept()
		if err != nil {
			fmt.Println("error :", err)
			os.Exit(1)
		}
		fmt.Println("connection 2 accept ")
		go copyConn(conn1, conn2, cChannel1, sigChannel)
		go copyConn(conn2, conn1, cChannel2, sigChannel)
	}
	// }

}

// func conListner(listener net.Listener, channelTx chan net.Conn, channelRx chan net.Conn, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for {
// 		conn1, err := listener.Accept()
// 		if err != nil {
// 			fmt.Println("error :", err)
// 			os.Exit(1)
// 		}
// 		fmt.Println("connection  accept ")
// 		channelTx <- conn1
// 		conn2 := <-channelRx
// 		go copyConn(conn1, conn2)
// 	}

// 	defer listener.Close()
// }

func copyConn(writer net.Conn, reader net.Conn, c chan int, sigChannel chan int) {
	// defer writer.Close()
	// defer reader.Close()

	select {
	case msg1 := <-sigChannel:
		fmt.Println("copyConn sigChannel:", msg1)
		return
	default:
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Println("io.Copy error:", err)
		} else {
			fmt.Println("copyConn connection closed")
		}
		reader.Close()
		c <- 0
		fmt.Println("copyConn connection before break")
		break
	}
	sigChannel <- 0
	fmt.Println("copyConn finished")
}

// func copyConn(writer, reader net.Conn) {
// 	defer writer.Close()
// 	defer reader.Close()

// 	_, err := io.Copy(writer, reader)
// 	if err != nil {
// 		fmt.Println("io.Copy error:", err)
// 	}

// 	fmt.Println("copyConn")
// }

// func tunnel(conn1 net.Conn, conn2 net.Conn) {

// }
