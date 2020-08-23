package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const serverPort = "3245"  //3245
const Err_Not_Valid_Token = "Token is not valid"
const Err_Use_Token_In_multiple_Device = "Token is used by another device"

// tokens = []string{"asdasdasd"}
var tokens = map[string]bool{}

type RequestData struct {
	Token string `json:"token"`
}

func (res *RequestData) UnmarshalJSON(buf []byte) {
	json.Unmarshal(buf, &res)
}

type ResponseData struct {
	Valid   int    `json:"valid"`
	Message string `json:"message"`
	PortNum int    `json:"portnum"`
}

func main() {
	tokens["123456"] = false
	tokens["654321"] = false
	// fmt.Println(tokens["654321"])
	// fmt.Println(tokens["6543213"])
	fmt.Println("Server Port:", serverPort)
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe(":"+serverPort, nil)

}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	// case "GET":
	// 	sendResponse(w, r)
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("requestHandler failed to ioutil.ReadAll(resp.Body):", err)
			os.Exit(1)
		}
		var data RequestData
		data.UnmarshalJSON(body)
		sendResponse(data, w, r)
	default:
		http.Error(w, "Sorry, only GET and POST methods are supported.", http.StatusNotFound)
		fmt.Println("Sorry, only GET and POST methods are supported.")
	}
}

func sendResponse(data RequestData, w http.ResponseWriter, r *http.Request) {
	var res ResponseData
	if _, ok := tokens[data.Token]; ok {
		// if val {
		// 	res.Valid = 0
		// 	res.Message = Err_Use_Token_In_multiple_Device
		// } else {
		res.Valid = 1
		res.Message = ""
		tokens[data.Token] = true
		res.PortNum = 4567
		// }
	} else {
		res.Valid = 0
		res.Message = Err_Not_Valid_Token
	}

	resMessage, _ := json.Marshal(res)
	fmt.Println("res :", string(resMessage))
	fmt.Fprintf(w, string(resMessage))

}
