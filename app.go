package main

type app struct {
	name    string
	version string
	receive bool
	send    bool
	file    string
}

func newApp() *app {
	return &app{
		name:    name,
		version: version,
		receive: false,
		send:    false,
		file:    "",
	}
}
