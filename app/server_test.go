package main

import (
	"fmt"
	"net"
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

// func TestMultipleRequestsInOneSession(t *testing.T) {

// 	conn, err := net.Dial("tcp", "localhost:4221")
// 	if err != nil {
// 		t.Fatalf("Failed to connect to server: %v", err)
// 	}
// 	conn.SetDeadline(time.Now().Add(10 * time.Second))

// 	request := "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"
// 	_, err = conn.Write([]byte(request))
// 	if err != nil {
// 		t.Fatalf("Failed to write to server: %v", err)
// 	}

// 	buffer := make([]byte, 4096)
// 	num, err := conn.Read(buffer)
// 	if err != nil {
// 		t.Fatalf("Failed to read from server: %v", err)
// 	}

// 	response := string(buffer[:num])
// 	fmt.Println("Response from server:", response)

//		//clear buffer
//		buffer = make([]byte, 4096)
//		request = "PUT /resource HTTP/1.1\r\nHost: localhost\r\n\r\n"
//		_, err = conn.Write([]byte(request))
//		if err != nil {
//			t.Fatalf("Failed to write to server: %v", err)
//		}
//		//read
//		num, err = conn.Read(buffer)
//		if err != nil {
//			t.Fatalf("Failed to read from server: %v", err)
//		}
//		//output
//		response = string(buffer[:num])
//		fmt.Println("Response from server:", response)
//		conn.Close()
//	}
func TestCloseAndOpenConn(t *testing.T) {
	//test
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
	conn, err := net.Dial("tcp", "localhost:4221")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	buffer := make([]byte, 4096)
	request := "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"
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
	_, err = conn.Write([]byte(request))
	if err != nil {
		t.Fatalf("Failed to write to server: %v", err)
	}
	request2 := "GET /whatsup HTTP/1.1\r\nHost: localhost\r\n\r\n"
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

	//clear buffer
	// buffer = make([]byte, 4096)
	// request = "PUT /resource HTTP/1.1\r\nHost: localhost\r\n\r\n"
	// _, err = conn.Write([]byte(request))
	// if err != nil {
	// 	t.Fatalf("Failed to write to server: %v", err)
	// }
	// //read
	// num, err = conn.Read(buffer)
	// if err != nil {
	// 	t.Fatalf("Failed to read from server: %v", err)
	// }
	// //output
	// response = string(buffer[:num])
	// fmt.Println("Response from server:", response)
	conn.Close()
	conn2.Close()
}

// TestParseTop tests the parseTop function.
// func TestParseTop(t *testing.T) {
// 	tests := []struct {
// 		topLine        string
// 		expectedMethod string
// 		expectError    bool
// 	}{
// 		{"GET / HTTP/1.1", GET, false},
// 		{"PUT /resource HTTP/1.1", PUT, false},
// 		{"INVALID / HTTP/1.1", "", true},
// 	}

// 	for _, test := range tests {
// 		method, _, _ := parseTop(test.topLine)
// 		if test.expectError {
// 			if method != "" {
// 				t.Errorf("Expected an error for input: %s", test.topLine)
// 			}
// 		} else {
// 			if method != test.expectedMethod {
// 				t.Errorf("parseTop(%s) = %s; want %s", test.topLine, method, test.expectedMethod)
// 			}
// 		}
// 	}
// }
