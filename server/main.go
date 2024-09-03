package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello, World!")

	apiServer := NewApiServer(":8080")
	err := apiServer.Run()
	if err != nil {
		log.Println(err)
	}
}
