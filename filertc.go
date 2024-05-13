package main

import (
	"log"

	"github.com/joho/godotenv"
)

const (
	version         = "1.0.0"
	name            = "filertc"
	senderBuffSize  = 16384
	bufferThreshold = 512 * 1024
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	setup()
}

func main() {
	app := newApp()
	if err := app.start(); err != nil {
		log.Fatal(err)
	}
}
