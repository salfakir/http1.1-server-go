package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

// TestIsValidMethod tests the isValidMethod function.
func TestIsValidMethod(t *testing.T) {
	tests := []struct {
		method   string
		expected bool
	}{
		{GET, true},
		{PUT, true},
		{POST, true},
		{DELETE, true},
		{"PATCH", false},
		{"OPTIONS", false},
	}

	for _, test := range tests {
		result := isValidMethod(test.method)
		if result != test.expected {
			t.Errorf("isValidMethod(%s) = %v; want %v", test.method, result, test.expected)
		}
	}
}

func TestCloseAndOpenConn(t *testing.T) {
	fmt.Println("TestCloseAndOpenConn")
	conn, err := net.Dial("tcp", "localhost:4221")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	conn.Close()

	conn, err = net.Dial("tcp", "localhost:4221")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	conn.Close()
}
func TestSendCloseAndOpenConn(t *testing.T) {
	fmt.Println("TestSendCloseAndOpenConn")
	conn, err := net.Dial("tcp", "localhost:4221")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	buffer := make([]byte, 4096)
	request := "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"
	fmt.Println("writting: " + request)
	_, err = conn.Write([]byte(request))
	if err != nil {
		t.Fatalf("Failed to write to server: %v", err)
	}
	//read
	num, err := conn.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read from server: %v", err)
	}
	//output
	response := string(buffer[:num])
	fmt.Println("Response from server:", response)

	conn.Close()
	conn, err = net.Dial("tcp", "localhost:4221")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	buffer = buffer[:0]

	request = "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"
	fmt.Println("writing: " + request)
	_, err = conn.Write([]byte(request))
	if err != nil {
		t.Fatalf("Failed to write to server: %v", err)
	}
	//read
	num, err = conn.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read from server: %v", err)
	}
	//output
	response = string(buffer[:num])
	fmt.Println("Response from server:", response)
	conn.Close()
}

func TestMultipleConnection(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:4221")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	conn2, err := net.Dial("tcp", "localhost:4221")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	conn2.SetDeadline(time.Now().Add(10 * time.Second))

	request := "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"
	fmt.Println("writing: " + request)
	_, err = conn.Write([]byte(request))
	if err != nil {
		t.Fatalf("Failed to write to server: %v", err)
	}
	request2 := "GET /whatsup HTTP/1.1\r\nHost: localhost\r\n\r\n"
	fmt.Println("writing: " + request2)
	_, err = conn2.Write([]byte(request2))
	if err != nil {
		t.Fatalf("Failed to write to server: %v", err)
	}

	buffer := make([]byte, 4096)
	num, err := conn.Read(buffer)
	if err != nil {
		t.Fatalf("Failed to read from server: %v", err)
	}

	response := string(buffer[:num])
	fmt.Println("Response from server for first request:", response)

	buffer2 := make([]byte, 4096)
	num2, err := conn2.Read(buffer2)
	if err != nil {
		t.Fatalf("Failed to read from server: %v", err)
	}
	response2 := string(buffer2[:num2])
	fmt.Println("Response from server for second request:", response2)

	conn.Close()
	conn2.Close()
}
func TestHelloGetRequest(t *testing.T) {
	// Create a client
	client := &http.Client{}

	// Build a request
	req, err := http.NewRequest("GET", "http://localhost:4221", nil)
	if err != nil {
		panic(err)
	}

	// Add some custom request headers
	req.Header.Set("User-Agent", "MyGoClient/1.0")
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// ---- Print Response Info ----
	// Status code
	fmt.Println("Status:", resp.Status)
	fmt.Println("StatusCode:", resp.StatusCode)

	// Response headers
	fmt.Println("\nResponse Headers:")
	for key, values := range resp.Header {
		for _, v := range values {
			fmt.Printf("%s: %s\n", key, v)
		}
	}

	// Response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nResponse Body:")
	fmt.Println(string(body))
}
func TestEchoGetRequest(t *testing.T) {
	// Create a client
	client := &http.Client{}
	randString := "printthisout!][()}+jli"
	uri := url.QueryEscape(randString)

	// Build a request
	req, err := http.NewRequest("GET", "http://localhost:4221/echo/"+uri, nil)
	if err != nil {
		panic(err)
	}

	// Add some custom request headers
	req.Header.Set("User-Agent", "MyGoClient/1.0")
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// ---- Print Response Info ----
	// Status code
	fmt.Println("Status:", resp.Status)
	fmt.Println("StatusCode:", resp.StatusCode)

	// Response headers
	fmt.Println("\nResponse Headers:")
	for key, values := range resp.Header {
		for _, v := range values {
			fmt.Printf("%s: %s\n", key, v)
		}
	}

	// Response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	sbody := string(body)
	//urldecode
	if sbody != randString {
		sbody, _ = url.QueryUnescape(sbody)
	}
	if sbody != randString {
		panic("Response body does not match the expected string, got " + sbody + ", expected " + randString)
	}
	fmt.Println("\nResponse Body:")
	fmt.Println(string(body))
}
func TestEchoHeaderGetRequest(t *testing.T) {
	// Create a client
	client := &http.Client{}

	// Build a request
	req, err := http.NewRequest("GET", "http://localhost:4221/user-agent", nil)
	if err != nil {
		panic(err)
	}

	// Add some custom request headers
	useragent := "MyGoClient/1.0"
	req.Header.Set("User-Agent", useragent)
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// ---- Print Response Info ----
	// Status code
	fmt.Println("Status:", resp.Status)
	fmt.Println("StatusCode:", resp.StatusCode)

	// Response headers
	fmt.Println("\nResponse Headers:")
	for key, values := range resp.Header {
		for _, v := range values {
			fmt.Printf("%s: %s\n", key, v)
		}
	}

	// Response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	sbody := string(body)
	//urldecode
	sbody, _ = url.QueryUnescape(sbody)
	if sbody != useragent {
		panic("Response body does not match the expected string, got " + sbody + ", expected " + useragent)
	}
	fmt.Println("\nResponse Body:")
	fmt.Println(string(body))
}
