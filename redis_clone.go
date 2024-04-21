package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var dataStore map[string]string

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		command := scanner.Text()

		parts := strings.Fields(command)
		fmt.Print(parts)

		if len(parts) < 2 {
			fmt.Fprintln(conn, "Invalid command")
			continue
		}

		switch {
		case parts[0] == "SET" || parts[0] == "set":
			if len(parts) != 3 {
				fmt.Fprintln(conn, "USE Syntax : SET key value")
				continue
			}
			key := parts[1]
			value := parts[2]
			dataStore[key] = value
			fmt.Fprintln(conn, "THE KEY VALUE IS ADDED")
		case parts[0] == "GET" || parts[0] == "get":
			if len(parts) != 2 {
				fmt.Fprintln(conn, "USE Syntax : GET key")
				continue
			}
			key := parts[1]
			value, ok := dataStore[key]
			if !ok {
				fmt.Fprintln(conn, "The key is not present in the dictionary")
				continue
			}
			fmt.Fprintln(conn, value)
		default:
			fmt.Fprintln(conn, "we have only SET AND GET METHOD")
		}
	}
}

func main() {
	dataStore = make(map[string]string)

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Redis clone listening on port 6379")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Print(conn)
		go handleConnection(conn)
	}
}
