package main

import "io"

type compHandler func()

type configuration struct {
	DPProvider   io.Reader   // The SDP reader
	SDPOutput    io.Writer   // The SDP writer
	OnCompletion compHandler // Handler to call on session completion
	STUN         string      // Custom STUN server
}

type commConfiguration struct {
	SDPProvider  io.Reader   // The SDP reader
	SDPOutput    io.Writer   // The SDP writer
	OnCompletion compHandler // Handler to call on session completion
	STUN         string      // Custom STUN server
}

type sendConfig struct {
	commConfiguration
	Stream io.Reader // The Stream to read from
}

type receiveConfig struct {
	commConfiguration
	Stream io.Writer // The Stream to write to
}
