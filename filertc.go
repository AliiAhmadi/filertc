package main

import (
	"log"

	"github.com/joho/godotenv"
)

const (
	version = "1.0.0"
	name    = "filertc"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	setup()
}

func main() {
	app := newApp()
	if err := app.start(); err != nil {
		log.Fatal(err)
	}
}
