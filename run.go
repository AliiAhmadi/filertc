package main

import (
	"errors"
	"flag"
	"os"

	"github.com/sirupsen/logrus"
)

func (a *app) start() error {
	a.parseFlags()
	a.validate()

	var err error

	if a.receive {
		err = a.receiveHandler()
	} else {
		err = a.sendHandler()
	}

	return err
}

func (a *app) parseFlags() {
	flag.BoolVar(&a.receive, "receive", false, "want to be Receiver")
	flag.BoolVar(&a.send, "send", false, "want to be Sender")
	flag.StringVar(&a.file, "file", "", "File want to send")
	flag.StringVar(&a.output, "output", "out", "Name of received file")
	flag.StringVar(&a.secret, "secret", "", "Secret of sender")

	flag.Parse()
}

func (a *app) validate() {
	if a.receive && a.send {
		logrus.Fatal("must specify one of receive or send")
	}

	if !(a.receive || a.send) {
		logrus.Fatal("must specify you want to be receiver or sender")
	}

	if a.send && a.file == "" {
		logrus.Fatal("file can not be empty")
	}

	if a.receive && a.secret == "" {
		logrus.Fatal("you should enter secret")
	}

	if a.send {
		_, err := os.Stat(a.file)
		if errors.Is(err, os.ErrNotExist) {
			logrus.Fatalf("file %s does not exist", a.file)
		}
	}
}
