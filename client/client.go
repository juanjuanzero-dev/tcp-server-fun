package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {

	var clientName string
	flag.StringVar(&clientName, "client", "Juan", "name of the calling client")
	flag.Parse()
	fmt.Println(clientName)
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	defer conn.Close()

	// send data to the server
	data := []byte(fmt.Sprintf("Hello Server! I'm %v", clientName))
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
