package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func contains(values []string, text string) bool {
	for _, value := range values {
		if value == text {
			return true
		}
	}
	return false
}

func serve(conn net.Conn, wordPtr *string) {
	request, err := http.ReadRequest(bufio.NewReader(conn))

	buf := make([]byte, 1024)

	if err != nil {
		fmt.Println("Error reading request. ", err.Error())
		return
	}

	path := request.URL.Path
	encoders := strings.Split(request.Header.Get("Accept-Encoding"), ", ")

	if strings.Contains(path, "/echo") {
		if len(encoders) > 1 {
			if !contains(encoders, "gzip") {
				conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n"))
			} else {
				conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: gzip\r\n\r\n"))
				conn.Close()
				return
			}
		} else {
			body := strings.Split(path, "/")[2]
			conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length:" + strconv.Itoa(len([]byte(body))) + "\r\n\r\n" + body))
		}
	} else if strings.Contains(path, "/files") {
		fileName := strings.Split(path, "/")[2]
		dat, err := os.ReadFile(*wordPtr + fileName)
		if err != nil {
			if request.Method == "POST" {
				n, _ := request.Body.Read(buf)
				os.WriteFile(*wordPtr+fileName, buf[:n], 0644)
				conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
				conn.Close()
				return
			}
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			conn.Close()
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

	wordPtr := flag.String("directory", "/tmp", "Directory to serve files from")
	flag.Parse()
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

		go serve(conn, wordPtr)
	}
}
