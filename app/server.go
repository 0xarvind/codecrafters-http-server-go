package main

import (
	"fmt"
	"strconv"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	//
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buf := make([]byte, 1024)
	conn.Read(buf)
	fmt.Println("Request received: " + string(buf))
	path := strings.Split(string(buf), " ")[1]
	if strings.Contains(path, "/echo") {
		body := strings.Split(path, "/echo")[2]
		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length:" + strconv.Itoa(len([]byte(body))) + "\r\n\r\n" + strings.Split(path, "/")[2]))
	} else if strings.HasPrefix(string(buf), "GET / HTTP/1.1") {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
	}
	conn.Close()
}
