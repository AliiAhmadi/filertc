package main

import "os"

func (a *app) sendHandler() error {
	f, err := os.Open(a.file)
	if err != nil {
		return err
	}
	defer f.Close()

	stun := os.Getenv("STUN_SERVER")

	err = parseSTUN(stun)
	if err != nil {
		return err
	}

	cf := sendConfig{
		Stream: f,
		commConfiguration: commConfiguration{
			OnCompletion: func() {},
			STUN:         stun,
		},
	}

	s := newSendSession(cf)
	s.a = a

	return s.start()
}
