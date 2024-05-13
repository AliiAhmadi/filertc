package main

import "os"

func (a *app) receiveHandler() error {
	f, err := os.OpenFile(a.file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	cf := receiveConfig{
		Stream: f,
	}

	s := newReceiveSession(&cf)

	return s.start()
}
