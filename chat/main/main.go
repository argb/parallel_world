package main

import (
	"log"
	"os"
	"parallel_world/chat/app"
)

func main() {
	f, err := os.OpenFile("/tmp/pworld.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	defer f.Close()

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)

	log.Println("This is a test log entry")

	myapp := app.NewApp()

	myapp.Run()
}
