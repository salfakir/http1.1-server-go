package main

import (
	"fmt"
	"strings"
	"time"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

// Define the allowed values as constants
const (
	GET    = "GET"
	PUT    = "PUT"
	POST   = "POST"
	DELETE = "DELETE"
)

type http_header struct {
	//define the struct for the headers
	name  string
	value string
}
type http_body struct {
	//define the struct for the body
	content string
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	startServer()
}
func startServer() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	//do while loop
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleParseRequest(conn)
	}
}

func handleParseRequest(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(200 * time.Second))
	if conn == nil {
		return
	}
	input := make([]byte, 4096)
	num, err := conn.Read(input)

	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
		return
	}

	//split the input into lines
	lines := string(input[:num])
	linearr := strings.Split(lines, "\r\n")
	http_top := linearr[0]
	http_method, http_path, http_version := parseTop(http_top, conn)

	//check if connection is closed
	if conn == nil {
		return
	}
	http_method = http_method + ""
	http_path = http_path + ""
	http_version = http_version + ""

	//remove the first line
	linearr = linearr[1:]
	headers, body := parseRequest(linearr, conn)
	if conn == nil {
		return
	}
	// headers = append(headers, http_header{name: "Connection", value: "close"})
	body.content = body.content + ""
	handleRequest(http_method, http_path, http_version, headers, body, conn)
}
func parseRequest(linearr []string, conn net.Conn) ([]http_header, http_body) {
	headers := []http_header{}
	body := http_body{}
	isCurrentlyHeaders := true
	emptylineCount := 0
	for _, line := range linearr {
		if line == "" {
			isCurrentlyHeaders = false
			emptylineCount++
			if emptylineCount == 2 {
				break
			}
			continue
		}
		if isCurrentlyHeaders {
			parts := strings.Split(line, ":")
			if handleHttpError(len(parts) == 2, "400 Bad Request", conn) {
				return nil, http_body{}
			}
			headers = append(headers,
				http_header{name: parts[0],
					value: strings.TrimSpace(parts[1])})
		} else {
			body.content = line
		}
	}
	return headers, body
}
func parseTop(http_top string, conn net.Conn) (string, string, string) {
	//parse the top line
	arr := strings.Split(http_top, " ")
	if handleHttpError(httpTopLength(arr), "400 Bad Request", conn) {
		return "", "", ""
	}
	http_method := arr[0]
	if handleHttpError(isValidMethod(http_method), "405 Method Not Allowed", conn) {
		return "", "", ""
	}
	http_path := arr[1]
	http_version := arr[2]
	if handleHttpError(checkVersion(http_version), "505 HTTP Version Not Supported", conn) {
		return "", "", ""
	}
	return http_method, http_path, http_version
}
func isValidMethod(http_method string) bool {
	//check if the method is valid
	switch http_method {
	case GET, PUT, POST, DELETE:
		return true
	}
	return false
}
func httpTopLength(arr []string) bool {
	if len(arr) != 3 {
		return false
	}
	return true
}
func handleHttpError(test bool, message string, conn net.Conn) bool {
	//assert the test
	if !test {
		_, err := conn.Write([]byte("HTTP/1.1 " + message + "\r\n"))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
		}
		conn.Close()
		return true
	}
	return false
}
func checkVersion(http_version string) bool {
	//check the version
	if http_version != "HTTP/1.1" {
		return false
	}
	return true
}
func handleRequest(http_method string, http_path string, http_version string,
	headers []http_header, body http_body, conn net.Conn) {
	//handle the request
	switch http_method + " " + http_path {
	case GET + " /":
		_, err := conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
		}
		fmt.Println("response for GET /")
	default:
		_, err := conn.Write([]byte("HTTP/1.1 404 Not Found my man\r\n\r\n"))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
		}
		fmt.Println("response for 404")
	}
}
