package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var dataStore map[string]expirationData
var mu sync.Mutex

type expirationData struct {
	value  string
	expiry time.Duration
}

func deleteExpiredKeys() {
	for {
		for key, expData := range dataStore {
			fmt.Println(key, expData)

			if expData.expiry != 0 {
				ticker := time.NewTicker(expData.expiry) // Tick every second
				for range ticker.C {
					fmt.Println(expData.expiry)
					mu.Lock()
					delete(dataStore, key)
					mu.Unlock()
					ticker.Stop()
					break
				}

			}
		}
	}
}

//}

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
			mu.Lock()
			dataStore[key] = expirationData{value: value}
			mu.Unlock()
			fmt.Fprintln(conn, "THE KEY VALUE IS ADDED")
		case command == "get":
			if len(parts) != 2 {
				fmt.Fprintln(conn, "USE Syntax : GET key")
				continue
			}
			key := parts[1]
			mu.Lock()
			expData, ok := dataStore[key]
			mu.Unlock()
			if !ok {
				fmt.Fprintln(conn, "The key is not present in the dictionary")
				continue
			}
			fmt.Fprintln(conn, expData.value)
		case command == "setex":
			if len(parts) != 4 {
				fmt.Fprintln(conn, "USE Syntax : SETEX key seconds value ")
				continue
			}
			key := parts[1]
			value := parts[3]
			expirySeconds, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Conversion failed:", err)
				return
			}

			expiryTime := time.Duration(expirySeconds) * time.Second

			mu.Lock()
			dataStore[key] = expirationData{value: value, expiry: expiryTime}
			mu.Unlock()
			fmt.Fprintf(conn, "THE KEY VALUE IS ADDED FOR ONLY %d Seconds\n", expirySeconds)
		default:
			fmt.Fprintln(conn, "we have only SET, GET, and SETEX methods")
		}
	}
}

func main() {
	dataStore = make(map[string]expirationData)

	go deleteExpiredKeys()

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
		fmt.Println(conn)
		go handleConnection(conn)
	}
}
