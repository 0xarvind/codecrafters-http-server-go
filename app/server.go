package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func serve(conn net.Conn) {
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		fmt.Println("Error reading request. ", err.Error())
		return
	}

	path := request.URL.Path
	fmt.Println("Path: ", path)
	if strings.Contains(path, "/echo") {
		body := strings.Split(path, "/")[2]
		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length:" + strconv.Itoa(len([]byte(body))) + "\r\n\r\n" + body))
	} else if strings.Contains(path, "/files") {
		fileName := strings.Split(path, "/")[2]
		dat, err := os.ReadFile("/tmp/" + fileName)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			return
		}
		fmt.Print(string(dat))
		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length:" + strconv.Itoa(len([]byte(dat))) + "\r\n\r\n" + string(dat)))
	} else if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if path == "/user-agent" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length:" + strconv.Itoa(len([]byte(request.UserAgent()))) + "\r\n\r\n" + request.UserAgent()))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
	conn.Close()
}

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
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go serve(conn)
	}
}
