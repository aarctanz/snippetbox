package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	connstr := os.Getenv("POSTGRESURL")
	if connstr == "" {
		log.Fatalln("Error loading postgres url")
	}
	db, err := NewDB(connstr)

	if err != nil {
		log.Fatalln(err)
	}

	apiServer := NewApiServer(":8080", db)
	err = apiServer.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
