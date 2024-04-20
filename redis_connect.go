package main

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

var key, value, live, check, gets string

// function is used to SET THE VALUE OF CERTAIN TIME.

func ttl_set(conn redis.Conn, err error) {
	_, err = conn.Do("SET", key, value, "EX", live)
	if err != nil {
		log.Fatal("Error setting value in Redis:", err)
	}
}

// function is used to GET the value
func getter(conn redis.Conn) {
	exists, err := redis.Bool(conn.Do("EXISTS", gets))
	if err != nil {
		panic(err)
	}
	if exists {
		val, err := redis.String(conn.Do("GET", gets))
		if err != nil {
			log.Fatal("Error getting value from Redis:", err)
		} else {
			fmt.Printf("Value of %v in Redis: %v\n", gets, val)
		}
	} else {
		fmt.Printf("Key %v does not exist in the database.\n", gets)

	}

}

// function is used to SET the value

func set(conn redis.Conn, err error) {

	_, err = conn.Do("SET", key, value)
	if err != nil {
		log.Fatal("Error setting value in Redis:", err)
	}
	fmt.Printf("Value of %v is updated \n", key)
}

func main() {
	// Connect to Redis server.
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal("Error connecting to Redis server:", err)
	}
	fmt.Println("Connected to Redis server")

	fmt.Printf("Enter your key : ")
	fmt.Scan(&key)
	fmt.Printf("Enter your value : ")
	fmt.Scan(&value)
	fmt.Printf("Do you need expiry options Y/N : ")
	fmt.Scan(&check)

	defer conn.Close()

	if check == "y" || check == "Y" {
		fmt.Printf("Set time to Live: ")
		fmt.Scan(&live)
		ttl_set(conn, err)

	} else {
		set(conn, err)
	}

	fmt.Printf("Enter your key to get : ")
	fmt.Scan(&gets)
	if check == "y" || check == "Y" {
		ttl, err := redis.Int(conn.Do("TTL", gets))
		if err != nil {
			panic(err)
		}
		if ttl > 0 {
			getter(conn)

		} else {
			fmt.Println("The value is dead ")
		}
	} else {
		getter(conn)
	}

}
