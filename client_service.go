package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

type Endpoint struct {
	Host string
	Port int
}

type RequestData struct {
	Token string `json:"token"`
}

type ResponseData struct {
	Valid   int    `json:"valid"`
	Message string `json:"message"`
	PortNum int    `json:"portnum"`
}

func (res *ResponseData) UnmarshalJSON(buf []byte) {
	json.Unmarshal(buf, &res)
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// const serverPort = "3245"
// var connection1 = &Endpoint{
// 	Host: "185.206.94.234",
// 	Port: 4570,
// }

var connection1 = &Endpoint{
	Host: "192.168.2.2",
	Port: 4570,
}

// var connection1 = &Endpoint{
// 	Host: "0.0.0.0",
// 	Port: 4570,
// }

// const serverURL = "http://0.0.0.0:3245"

const serverURL = "http://192.168.2.2:3245"

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Please Provide token")
		os.Exit(1)
	}
	token := os.Args[1]

	if !checkToken(token) {
		os.Exit(0)
	}

	connection2 := &Endpoint{
		Host: "localhost",
		Port: 22,
	}

	// connection1 := &Endpoint{
	// 	Host: "",
	// 	Port: 4570,
	// }

	// listener, err := net.Listen("tcp", connection1.String())
	// if err != nil {
	// 	fmt.Println("connection1 error :", err)
	// 	os.Exit(1)
	// }

	conn1, err := net.Dial("tcp", connection1.String())
	if err != nil {
		fmt.Println("connection 1 error: ", err)
		os.Exit(1)
	}

	fmt.Println("Connection 1 connected")

	conn2, err := net.Dial("tcp", connection2.String())
	if err != nil {
		fmt.Printf("connection 2 error: %s", err)
		os.Exit(1)
	}

	// defer listener.Close()
	defer conn1.Close()
	defer conn2.Close()

	go copyConn(conn1, conn2)
	go copyConn(conn2, conn1)
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

func checkToken(token string) bool {
	var req RequestData
	req.Token = token
	reqMessage, _ := json.Marshal(req)
	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(reqMessage))
	if err != nil {
		fmt.Println("requestHandler failed to http.Post:", err)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("requestHandler failed to ioutil.ReadAll(resp.Body):", err)
		os.Exit(1)
	}
	var data ResponseData
	data.UnmarshalJSON(body)

	if data.Valid == 1 {
		connection1.Port = data.PortNum
		return true
	}

	fmt.Println(data.Message)
	return false
}
