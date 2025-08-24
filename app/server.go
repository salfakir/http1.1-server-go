package main

import (
	"fmt"
	"io"
	"regexp"
	"runtime"
	"strconv"
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
type http_error struct {
	code     int
	message  string
	location *error_location
}
type error_location struct {
	file string
	line int
	ok   bool
}
type http_request struct {
	http_method  string
	http_path    string
	http_version string
	headers      []http_header
	body         http_body
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
		if err.Error() == "EOF" {
			// Connection closed by client
			return
		}
		fmt.Println("Error reading data: ", err.Error())
		return
	}

	//split the input into lines
	lines := string(input[:num])
	linearr := strings.Split(lines, "\r\n")
	http_top := linearr[0]
	http_method, http_path, http_version, http_error := parseTop(http_top)
	if http_error.code != 0 {
		handleHttpError(http_error, conn)
		return
	}

	http_method = http_method + ""
	http_path = http_path + ""
	http_version = http_version + ""

	//remove the first line
	linearr = linearr[1:]
	headers, body, http_error := parseRequest(linearr)
	if http_error.code != 0 {
		handleHttpError(http_error, conn)
		return
	}
	// headers = append(headers, http_header{name: "Connection", value: "close"})
	body.content = body.content + ""
	req := http_request{
		http_method, http_path, http_version, headers, body,
	}
	handleRequest(req, conn)
}
func parseRequest(linearr []string) ([]http_header, http_body, http_error) {
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
			if len(parts) < 2 {
				_, file, line, ok := runtime.Caller(1)
				return nil, http_body{}, http_error{400, "Bad Request", &error_location{file, line, ok}}
			}
			headervalue := strings.Join(parts[1:], ":")
			headers = append(headers,
				http_header{name: parts[0],
					value: strings.TrimSpace(headervalue)})
		} else {
			body.content = line
		}
	}
	return headers, body, http_error{0, "", nil}
}
func parseTop(http_top string) (string, string, string, http_error) {
	//parse the top line
	arr := strings.Split(http_top, " ")
	if !httpTopLength(arr) {
		_, file, line, ok := runtime.Caller(1)
		return "", "", "", http_error{400, "Bad Request", &error_location{file, line, ok}}
	}
	http_method := arr[0]
	if !isValidMethod(http_method) {
		_, file, line, ok := runtime.Caller(1)
		return "", "", "", http_error{405, "Method Not Allowed", &error_location{file, line, ok}}
	}
	http_path := arr[1]
	http_version := arr[2]
	if !checkVersion(http_version) {
		_, file, line, ok := runtime.Caller(1)
		return "", "", "", http_error{505, "HTTP Version Not Supported", &error_location{file, line, ok}}
	}
	return http_method, http_path, http_version, http_error{0, "", nil}
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
	if len(arr) == 3 {
		return true
	}
	return false
}

func handleHttpError(he http_error, conn net.Conn) {
	message := he.message
	code := he.code
	line := ""
	file := ""
	if he.location != nil && he.location.ok {
		line = strconv.Itoa(he.location.line)
		file = he.location.file
	}
	fmt.Println("[" + file + ":" + line + "] Error parsing request(" + strconv.Itoa(he.code) + "): " + he.message)

	_, err := conn.Write([]byte("HTTP/1.1 " + strconv.Itoa(code) + " " + message + "\r\n\r\n"))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
	conn.Close()
}
func checkVersion(http_version string) bool {
	//check the version
	if http_version != "HTTP/1.1" {
		return false
	}
	return true
}
func handleRequest(req http_request, conn net.Conn) {
	//debug
	fmt.Println("Handling request: ", req.http_method, req.http_path, req.http_version)
	fmt.Println("Headers: ")
	for _, header := range req.headers {
		fmt.Println("  ", header.name, ":", header.value)
	}
	switch req.http_method {
	case GET:
		handleGet(req, conn)
	}
}
func handleGet(req http_request, conn net.Conn) {
	if req.http_path == "/index.html" || req.http_path == "/index.htm" || req.http_path == "/" {
		handleResponse("HTTP/1.1 200 OK",
			[]http_header{http_header{name: "Content-Type", value: "text/html"}},
			http_body{content: "<html><body><h1>Welcome to the index page!</h1></body></html>"},
			conn)
		fmt.Println("response for GET /index.html")
		return
	} else if regexp.MustCompile(`^/echo/[a-zA-Z0-9_\-%\(\)';:\+\*\$=\[\]]+$`).MatchString(req.http_path) {
		parts := strings.Split(req.http_path, "/")
		echo := parts[2]
		length := len(echo)
		handleResponse("HTTP/1.1 200 OK",
			[]http_header{
				http_header{name: "Content-Type", value: "text/plain"},
				http_header{name: "Content-Length", value: strconv.Itoa(length)},
			},
			http_body{content: echo},
			conn,
		)

		return
	} else if req.http_path == "/user-agent" {
		echo := ""
		for _, header := range req.headers {
			if strings.ToLower(header.name) == "user-agent" {
				echo = header.value
				break
			}
		}
		handleResponse("HTTP/1.1 200 OK",
			[]http_header{
				http_header{name: "Content-Type", value: "text/plain"},
				http_header{name: "Content-Length", value: strconv.Itoa(len(echo))},
			},
			http_body{content: echo},
			conn,
		)
		return
	} else if regexp.MustCompile(`^/files/[a-zA-Z0-9_\-]+$`).MatchString(req.http_path) {
		file := strings.TrimPrefix(req.http_path, "/files/")
		file = strings.Trim(file, "/")
		path := "/tmp/" + file
		if _, err := os.Stat(path); os.IsNotExist(err) {
			handleNotFound(conn)
			return
		}
		f, err := os.Open(path)
		if err != nil {
			handleInternalError(conn)
			return
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		sc := string(content)
		if err != nil {
			handleInternalError(conn)
			return
		}
		length := len(sc)
		handleResponse("HTTP/1.1 200 OK",
			[]http_header{
				http_header{name: "Content-Type", value: "application/octet-stream"},
				http_header{name: "Content-Length", value: strconv.Itoa(length)},
			},
			http_body{content: sc},
			conn,
		)
		return
	} else {
		handleNotFound(conn)
	}
}
func handleResponse(top string, headers []http_header, body http_body, conn net.Conn) {
	//debug
	fmt.Println("Handling response: ", top)
	fmt.Println("Headers: ")
	for _, header := range headers {
		fmt.Println("  ", header.name, ":", header.value)
	}
	response := top + "\r\n"
	for _, header := range headers {
		if len(header.name) == 0 || len(header.value) == 0 {
			continue
		}
		response += header.name + ": " + header.value + "\r\n"
	}
	trimbody := strings.TrimRight(body.content, "\r\n")
	response += "\r\n" + trimbody + "\r\n\r\n"
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
}
func handleNotFound(conn net.Conn) {
	handleResponse("HTTP/1.1 404 Not Found",
		[]http_header{},
		http_body{content: ""},
		conn)
}
func handleInternalError(conn net.Conn) {
	handleResponse("HTTP/1.1 500 Internal Server Error",
		[]http_header{},
		http_body{content: ""},
		conn)
}
