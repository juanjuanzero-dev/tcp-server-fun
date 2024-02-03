package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Listen for incoming connections
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")
	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error", err)
			continue
		}
		// handle client connection in a go routine
		go handleClient(conn)
	}
}
func handleClient(conn net.Conn) {
	defer conn.Close()

	// Create a buffer to read data into
	buffer := make([]byte, 1024)
	for { // a while(true) loop
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil && n == 0 {
			return // nothing more is read
		}
		if err != nil {
			fmt.Printf("Error: %v, Bytes read: %v\n", err, n)
			return
		}

		// Process and use the data
		fmt.Printf("Received: %s\n", buffer[:n]) // reading from 0 to n not inclusive
		writeLogV2(buffer[:n])
	}
}

// writes to a file, but updates the file everytime
func writeLog(message []byte) {
	//create the file
	logFile, err := os.Create("logfile.txt")
	if err != nil {
		fmt.Println("error:", err)
	}
	bytesWritten, err := logFile.Write(message)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("Bytes:", bytesWritten)
}

// writes to a file, but does not truncate, appends if it already exists.
func writeLogV2(message []byte) {
	logFile, err := os.OpenFile("logfile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("error:", err)
	}

	bytesWritten, err := logFile.Write(message)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("Bytes:", bytesWritten)
	logFile.Close()
}
