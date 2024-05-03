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

var mu sync.Mutex

type NewQueue struct {
	items []string
}

type expirationData struct {
	value  string
	expiry time.Time
}

var dataStore map[string]expirationData

var queue map[string]NewQueue
var ch = make(chan string)

func deleteExpiredKeys() {
	for range time.Tick(10 * time.Second) {
		mu.Lock()
		for key, expData := range dataStore {
			if !expData.expiry.IsZero() && time.Now().After(expData.expiry) {
				delete(dataStore, key)
			}
		}
		mu.Unlock()
	}
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
			fmt.Println(key)
			fmt.Println(value)

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
			fmt.Println(key)
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
			fmt.Print(time.Duration(expirySeconds) * time.Second)
			expiryTime := time.Now().Add(time.Duration(expirySeconds) * time.Second)
			fmt.Print(expiryTime)

			mu.Lock()
			dataStore[key] = expirationData{value: value, expiry: expiryTime}
			mu.Unlock()
			fmt.Fprintf(conn, "THE KEY VALUE IS ADDED FOR ONLY %d Seconds\n", expirySeconds)
		// lpush list name
		case command == "lpush":
			if len(parts) < 3 {
				fmt.Fprintln(conn, "USE Syntax: LPUSH key value")
				continue
			}
			key := parts[1]
			value := parts[2]

			mu.Lock()
			ed, ok := queue[key]
			if !ok {
				ed = NewQueue{items: []string{value}}
			} else {
				ed.items = append(ed.items, value)
			}
			queue[key] = ed
			fmt.Println(ed.items)

			mu.Unlock()
			fmt.Fprintln(conn, "ELEMENT ADDED TO LIST")
		// blpop list 10
		case command == "blpop":
			if len(parts) != 3 {
				fmt.Fprintln(conn, "USE Syntax: BLPOP key timeout")
				continue
			}
			key := parts[1]
			timeout, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Fprintf(conn, "Invalid timeout value: %s\n", parts[2])
				continue
			}

			mu.Lock()

			list, exists := queue[key]
			if !exists || len(list.items) == 0 {
				expiry := time.Now().Add(time.Duration(timeout) * time.Second)
				go func() {
					for {
						mu.Lock()
						list, exists := queue[key]
						mu.Unlock()

						if exists && len(list.items) > 0 {
							fmt.Println("hi")

							ch <- list.items[0]
							break
						}
						if time.Now().After(expiry) && timeout != 0 {
							fmt.Println("hello")
							ch <- "nil"
							break
						}
					}
				}()
			} else {
				// Pop the first item from the list
				item := list.items[0]
				list.items = list.items[1:]
				queue[key] = list
				fmt.Fprintln(conn, item)
				mu.Unlock()
				continue
			}
			mu.Unlock()
			select {

			case item := <-ch:
				if item != "nil" {
					queue[key] = NewQueue{items: []string{}}

				}
				fmt.Fprintln(conn, item)
			}
		// lrange mylist 0 -1
		case command == "lrange":
			if len(parts) != 4 {
				fmt.Fprintln(conn, "USE Syntax: LRANGE key start stop")
				continue
			}
			key := parts[1]
			start, err := strconv.Atoi(parts[2])
			stop, err := strconv.Atoi(parts[3])
			if err != nil {
				fmt.Fprintln(conn, "You have given wrong start or stop value")

				continue
			}

			list, exists := queue[key]
			total_len := len(list.items)
			//if start > total_len || stop > total_len {
			//	fmt.Fprintln(conn, "Queue doesn't have that much element")
			//	continue

			//}
			if !exists || total_len == 0 {
				fmt.Fprintln(conn, " The queue does not exist")
			} else {
				fmt.Fprintln(conn, list.items[start:stop])
			}

		default:
			fmt.Fprintln(conn, "we have only SET, GET, and SETEX methods")

		}
	}
}

func main() {
	dataStore = make(map[string]expirationData)
	queue = make(map[string]NewQueue)
	//ch = make(chan string)

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
