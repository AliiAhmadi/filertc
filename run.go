package main

import "flag"

func (a *app) start() error {
	a.parseFlags()

	return nil
}

func (a *app) parseFlags() {
	flag.BoolVar(&a.receive, "receive", false, "Receiver")
	flag.BoolVar(&a.send, "send", false, "Sender")
	flag.StringVar(&a.file, "file", "", "File want to send")

	flag.Parse()
}
