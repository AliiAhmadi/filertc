package main

import (
	"log"
)

const (
	version         = "1.0.0"
	name            = "filertc"
	senderBuffSize  = 16384
	bufferThreshold = 512 * 1024
	stun            = "stun.l.google.com:19302"
	local_file      = "local_sdp.txt"
	remote_file     = "remote_sdp.txt"
)

func init() {
	setup()
}

func main() {
	app := NewApp()
	if err := app.start(); err != nil {
		log.Fatal(err)
	}
}
