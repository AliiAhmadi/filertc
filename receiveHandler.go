package main

import "os"

func (a *app) receiveHandler() error {
	f, err := os.OpenFile(a.output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = parseSTUN(stun)
	if err != nil {
		return err
	}

	cf := receiveConfig{
		Stream: f,
		commConfiguration: commConfiguration{
			OnCompletion: func() {},
			STUN:         stun,
		},
	}

	s := newReceiveSession(cf)
	s.a = a

	return s.start()
}
