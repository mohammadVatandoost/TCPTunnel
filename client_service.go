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
	"strconv"
	"sync"
	"time"
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

var connection2 = &Endpoint{
	Host: "127.0.0.1",
	Port: 4530,
}

var connectionType = "tcp"

// var serverIp = "192.168.100.7"
var serverIp = "0.0.0.0"

// var connection1 = &Endpoint{
// 	Host: "0.0.0.0",
// 	Port: 4570,
// }

const serverURL = "http://0.0.0.0:3245"

// const serverURL = "http://192.168.2.2:3245"

func main() {

	if len(os.Args) != 4 {
		fmt.Println("Please Provide token and Connection info")
		os.Exit(1)
	}
	token := os.Args[1]
	connectionType = os.Args[2]
	localPortNum, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Please Provice valid port number, err:", err)
		os.Exit(0)
	}
	if !checkToken(token) {
		os.Exit(0)
	}

	connection2 := &Endpoint{
		Host: "127.0.0.1",
		Port: localPortNum,
	}

	var wg sync.WaitGroup
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

	fmt.Println("Connection 2 connected")

	fmt.Println("Tunnely")
	fmt.Println("Connection type: ", connectionType)
	fmt.Println("Forwarding localhost:", localPortNum, " to ", serverIp, ":", connection1.Port)

	defer conn1.Close()
	defer conn2.Close()
	if connectionType == "tcp" {
		wg.Add(1)
		go copyConn(conn1, conn2, &wg)
		wg.Add(1)
		go copyConn(conn2, conn1, &wg)
		wg.Wait()
		for range time.Tick(time.Minute * 1) {
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
			fmt.Println("Connection 2 connected")
			defer conn1.Close()
			defer conn2.Close()
			wg.Add(1)
			go copyConn(conn1, conn2, &wg)
			wg.Add(1)
			go copyConn(conn2, conn1, &wg)
			wg.Wait()
		}
	} else if connectionType == "http" {
		for {
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
			wg.Add(1)
			go copyConn(conn1, conn2, &wg)
			wg.Add(1)
			go copyConn(conn2, conn1, &wg)
			wg.Wait()
		}

		defer conn1.Close()
		defer conn2.Close()
	}

}

// fmt.Println("Server Port:", localPortNum)
// 		http.HandleFunc("/", requestHandler)
// 		http.ListenAndServe(":"+strconv.Itoa(connection1.Port), nil)

// func requestHandler(w http.ResponseWriter, req *http.Request) {
// 	fmt.Println("requestHandler")
// 	// we need to buffer the body if we want to read it here and send it
// 	// in the request.
// 	body, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// you can reassign the body if you need to parse it as multipart
// 	req.Body = ioutil.NopCloser(bytes.NewReader(body))

// 	// create a new url from the raw RequestURI sent by the client
// 	// url := fmt.Sprintf("%s://%s%s", proxyScheme, proxyHost, req.RequestURI)
// 	url := "http://" + connection2.Host + string(":") + strconv.Itoa(connection2.Port)
// 	proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))

// 	// We may want to filter some headers, otherwise we could just use a shallow copy
// 	// proxyReq.Header = req.Header
// 	proxyReq.Header = make(http.Header)
// 	for h, val := range req.Header {
// 		proxyReq.Header[h] = val
// 	}
// 	client := &http.Client{}
// 	resp, err := client.Do(proxyReq)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadGateway)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
// 	w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
// 	io.Copy(w, resp.Body)
// 	resp.Body.Close()
// }

// func pipeReq(rw http.ResponseWriter, req *http.Request) {
//     resp, err := http.Get(".....")
//     if err != nil{
//         //handle the error
//         return
//     }
//     rw.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
//     rw.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
//     io.Copy(rw, resp.Body)
//     resp.Body.Close()

// }

func copyConn(writer, reader net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer writer.Close()
	defer reader.Close()

	_, err := io.Copy(writer, reader)
	if err != nil {
		fmt.Println("io.Copy error:", err)
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
		connection1.Host = serverIp
		connection1.Port = data.PortNum
		return true
	}

	fmt.Println(data.Message)
	return false
}
