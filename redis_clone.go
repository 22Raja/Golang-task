package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var key, value string
var ttl int
var dataStore map[string]string
var command string

func kill(key string) {
	time.Sleep(time.Second * time.Duration(ttl))
	delete(dataStore, key)

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		command := scanner.Text()

		parts := strings.Fields(command)

		if len(parts) < 2 {
			fmt.Fprintln(conn, "Invalid command")
			continue
		}
		command = strings.ToLower(parts[0])

		switch {
		case command == "set":
			if len(parts) != 3 {
				fmt.Fprintln(conn, "USE Syntax : SET key value")
				continue
			}
			key := parts[1]
			value := parts[2]
			dataStore[key] = value
			fmt.Fprintln(conn, "THE KEY VALUE IS ADDED")
		case command == "get":
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
		//setex name 10 raja
		case command == "setex":
			if len(parts) != 4 {
				fmt.Fprintln(conn, "USE Syntax : SETEX key seconds value ")
				continue
			}
			value = parts[3]
			key = parts[1]
			dataStore[key] = value
			num, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Conversion failed:", err)
				return
			}
			ttl = num
			go kill(key)
			fmt.Fprintln(conn, "THE KEY VALUE IS ADDED FOR ONLY ", num, "Seconds")

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

