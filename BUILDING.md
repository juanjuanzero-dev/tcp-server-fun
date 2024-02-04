# Building a TCP/IP Server and CLI Uploader

Goal: Build a TCP IP service and cli uploader as a golang project. Just for fun

## What is TCP/IP?

Its a suite of protocols and standards that govern how data is transmitted and received across networks. It is hierarchical in nature, here are the layers (like Ogres)

- Application Layer: includes application specific protocols like HTTP, FTP and SMTP (emails) where day is consumed by end-user applications
- Transport Layer: this managers end-to-end communication. It includes 2 protocols:
  - TCP: reliable, ordered and error-checked data transmission between two devices
  - UDP: a lightweight, connectionless protocol ideal for real-time apps where low latency is crucial, like streaming or online gaming
- Internet Layer: where addressing and routing data packets across interconnected networks happen.
- Link layer: handles physical transmission of data packets over specific network medium, it manages the data link between devices on the same network

## Data Concepts

- Buffering: this involves temporarily storing received data in memory until it can be processed. Buffers help manage the rate at which data is read from and written to a connection. Buffer sizes can be adjusted to optimize performance.
- Streaming: the process of continuously sending or receiving data without waiting fro the entire data set to be available.

## What is a socket?

- Socket Address: like an IP and port pair, what identifies a specific endpoint in a network
- Server Sockets: listeners from clients
- Client Sockets: use to initiate connects to servers
- Socket Communication: enables bidirectional communication allowing data to be sent and received between connected devices.

## Making a server in golang

A very basic server and client. The Server listens for connections and handles the incomming request.

When handling TCP connections a few concepts occur:

- Data Chuncking: the data to be transmitted is divided into smaller units called segments
- Segmenting and Sequence Numbers: each segment is assigned a sequence number to ensure that the receiver can reassembel them in the correct order.
- Transmission: TCP guarantees that the segments are received reliably and in the correct order
- Ack: The receiver acks the receipt of segments back to the sender, if the sender doesnt receive an ack within a time frame, it retransmits the segment

### The Very First Cut

For the first cut we'll create a server that listens on a port and prints out the message from the client.

Here is a very simplest for of a TCP/IP server

```go
package main

import (
	"fmt"
	"net"
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
	}
}
```

Here is the client code in another project

```go
package main

import (
	"fmt"
	"net"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	defer conn.Close()

	// send data to the server
	data := []byte("Hello, server")
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

```

## Do something else with the data

- [x] Change the code so that we do something with the data, i.e. save to a log file for incomming messages.

- I learned more about the file system here, in creating files from scratch. It was nice and as it was simple. os.Create creates or truncates the file, where as, os.OpenFile can open a file or create it if it does not exist.

- I also learned about the numeric notation of on the file system. and `0666` means to have read and write permissions on the file as the file mode. In with the pipe characters is a flag that the file gets opened with. First is the O_APPEND flag to append to an existing file, if that is unsuccessful we use the O_CREATE flag to write to the file.

Here is the server code that takes in the files.

```go

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
```

- [x] Add a param in the client code that takes in a flag for the client name, this will also get logged.

- to implement this we'll need to pass in flags using the `flag` package in the standard library

- This took me a while to learn, mainly because im an junior trying to get better, its part of the journey eh? So I added code to create a flag. That gets parsed, Here is the code:

```go
	var clientName string
	flag.StringVar(&clientName, "client", "Juan", "name of the calling client") // wire up the flag
	flag.Parse() // Parse all flags
	fmt.Println(clientName) // print it out, for checking
```

What this does is very simple, it creats a string variable `clientName`, then we map the pointer to that so the value of the incomming flag called `client`, which defaults to `Juan` is the name of the calling client. Now all we have to do is, where the clientName will be `zero`:

```shell
go run client.go -client zero
```

- [x] make the file name of the file a param so you can have separate logs for each client.

- To implement this i'd like to add some structure to what is going on, so im electing to encode the message in a better way. This the server doesn't have to guess. Let's use JSON. The client will send JSON message and the server will pass it.
- we create a new struct call message that will house the clientName, and the message. Make sure to export the values and add struct tags (perhaps not needed) but was useful in realizing that the properties needed to be exported, what was happening before was we were sending out empty objects after Marshalling to JSON. Likely because it was not exported.
- on the server side we Unmarshall that JSON, and make our own files using the client name, we used the filepath package in the standard library to work with file paths. Also i got gopls to work in nvim

## Evolving to Handle Multiple Concurrent TCP Connections

- Change so that we are handling multiple concurrent TCP connections, implement a connection pool, and timeouts

## Evolving to Handle Errors and Handle Graceful Shutdown

- add a gracefull shut down

## Security, Implement Secure Connections

## Resource

https://okanexe.medium.com/the-complete-guide-to-tcp-ip-connections-in-golang-1216dae27b5a
