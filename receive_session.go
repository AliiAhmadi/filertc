package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pion/webrtc/v4"
	"github.com/sirupsen/logrus"
)

type receiveSession struct {
	sess        inSession
	stream      io.Writer
	msgChannel  chan webrtc.DataChannelMessage
	initialized bool
}

func newReceiveSession(c *receiveConfig) *receiveSession {
	s := &receiveSession{
		sess: inSession{
			sdpInput:     c.SDPProvider,
			sdpOutput:    c.SDPOutput,
			Done:         make(chan struct{}),
			NetworkStats: newStats(),
			stunServers:  []string{fmt.Sprintf("stun:%s", os.Getenv("STUN_SERVER"))},
		},
		stream:      c.Stream,
		msgChannel:  make(chan webrtc.DataChannelMessage, 4096*2),
		initialized: false,
	}

	if s.sess.sdpInput == nil {
		s.sess.sdpInput = os.Stdin
	}

	if s.sess.sdpOutput == nil {
		s.sess.sdpOutput = os.Stdout
	}

	return s
}

func (s *receiveSession) start() error {
	if err := s.initialize(); err != nil {
		return err
	}

	s.receiveData()
	s.sess.onCompletion()
	return nil
}

func (s *receiveSession) initialize() error {
	if s.initialized {
		return nil
	}

	if err := s.sess.createConnection(s.onConnectionStateChange()); err != nil {
		logrus.Errorln(err)
		return err
	}

	s.createDataHandler()
	if err := s.sess.readSDP(); err != nil {
		logrus.Errorln(err)
		return err
	}

	if err := s.sess.createAnswer(); err != nil {
		logrus.Errorln(err)
		return err
	}

	s.initialized = true
	return nil
}
