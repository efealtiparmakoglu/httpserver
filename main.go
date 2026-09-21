package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const version = "1.0.0"

// HTTP/1.1 server using raw net.Listener — no net/http. Routes are exact
// prefix matches. Content is generated per handler.

type Handler func(method, path, body string) (status int, contentType string, bodyOut string)

func main() {
	port := 8080
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-p" || os.Args[i] == "--port" {
			i++
			if i < len(os.Args) {
				port, _ = strconv.Atoi(os.Args[i])
			}
		}
		if os.Args[i] == "--help" || os.Args[i] == "-h" {
			fmt.Println("httpserver", version, "- raw HTTP/1.1 server, zero dependencies")
			fmt.Println("USAGE: httpserver [-p PORT]")
			os.Exit(0)
		}
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "httpserver: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("listening on :%d\n", port)

	handler := func(method, path, _ string) (int, string, string) {
		switch {
		case path == "/":
			return 200, "text/html", "<h1>httpserver " + version + "</h1><p>zero dependencies</p>"
		case path == "/health":
			return 200, "application/json", `{"status":"ok"}`
		case path == "/echo":
			return 200, "text/plain", method + " " + path
		case strings.HasPrefix(path, "/hello/"):
			return 200, "text/plain", "hello, " + strings.TrimPrefix(path, "/hello/")
		default:
			return 404, "application/json", `{"error":"not found"}`
		}
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "accept: %v\n", err)
			continue
		}
		go handleConn(conn, handler)
	}
}

func handleConn(conn net.Conn, handler Handler) {
	defer conn.Close()
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		return
	}
	request := string(buf[:n])

	// parse request line
	lines := strings.SplitN(request, "\r\n", 2)
	if len(lines) == 0 {
		return
	}
	parts := strings.Fields(lines[0])
	if len(parts) < 3 {
		return
	}
	method, path := parts[0], parts[1]

	// extract body if Content-Length present
	body := ""
	for _, line := range strings.Split(request, "\r\n\r\n")[1:] {
		body = line
	}
	for _, line := range strings.Split(lines[0], "\r\n")[1:] {
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "content-length:") {
			if cl, e := strconv.Atoi(strings.TrimSpace(line[15:])); e == nil && cl > 0 && cl <= len(body) {
				body = body[:cl]
			}
		}
	}

	status, contentType, bodyOut := handler(method, path, body)
	statusText := map[int]string{200: "OK", 404: "Not Found"}[status]
	if statusText == "" {
		statusText = "OK"
	}
	response := fmt.Sprintf(
		"HTTP/1.1 %d %s\r\nContent-Type: %s\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s",
		status, statusText, contentType, len(bodyOut), bodyOut)
	conn.Write([]byte(response))
}
