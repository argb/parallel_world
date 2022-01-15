package main

import (
	"log"
	"os"
	app2 "parallel_world/chat/app"
)

func main() {
	f, err := os.OpenFile("/tmp/pworld.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	defer f.Close()

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)

	log.Println("This is a test log entry")

	app := app2.NewApp()

	app.Run()
}
