package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	ClientName string `json:"clientName"`
	Message    string `json:"message"`
}

func main() {

	var clientName string
	flag.StringVar(&clientName, "client", "Juan", "name of the calling client")
	flag.Parse()
	fmt.Println(clientName)
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error,", err)
		return
	}
	defer conn.Close()

	// send data to the server
	var msg = Message{
		ClientName: clientName,
		Message:    fmt.Sprintf("Hey there server! this is %v", clientName),
	}
	fmt.Println(msg)

	data, err := json.Marshal(msg) // convert to JSON
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
