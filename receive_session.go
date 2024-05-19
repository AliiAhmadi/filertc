package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pion/webrtc/v2"
	"github.com/sirupsen/logrus"
)

type receiveSession struct {
	sess        inSession
	stream      io.Writer
	msgChannel  chan webrtc.DataChannelMessage
	initialized bool
	a           *app
}

func newReceiveSession(c receiveConfig) *receiveSession {
	s := &receiveSession{
		sess: inSession{
			sdpInput:     c.SDPProvider,
			sdpOutput:    c.SDPOutput,
			Done:         make(chan struct{}),
			NetworkStats: newStats(),
			stunServers:  []string{fmt.Sprintf("stun:%s", stun)},
			onCompletion: func() {},
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

func (s *receiveSession) receiveData() {
	logrus.Infoln("start receiving data ...")
	defer logrus.Infoln("done receiving ...")

	for {
		select {
		case <-s.sess.Done:
			s.sess.NetworkStats.stop()
			return
		case msg := <-s.msgChannel:
			n, err := s.stream.Write(msg.Data)
			if err != nil {
				logrus.Errorln(err)
			} else {
				speed := s.sess.NetworkStats.bandwidth()
				fmt.Printf("%.2f MB/s\r", speed)
				s.sess.NetworkStats.addBytes(uint64(n))
			}
		}
	}
}

func (s *receiveSession) onConnectionStateChange() func(webrtc.ICEConnectionState) {
	return func(i webrtc.ICEConnectionState) {
		logrus.Infof("ICE Connection State has changed: %s\n", i.String())
	}
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

func (s *receiveSession) createDataHandler() {
	s.sess.onDataChannel(func(d *webrtc.DataChannel) {
		logrus.Debugf("New DataChannel %s %d\n", d.Label(), d.ID())
		s.sess.NetworkStats.start()
		d.OnMessage(s.onMessage())
		d.OnClose(s.onClose())
	})
}

func (s *inSession) onDataChannel(f func(*webrtc.DataChannel)) {
	s.peerConnection.OnDataChannel(f)
}

func (s *inSession) createAnswer() error {
	ans, err := s.peerConnection.CreateAnswer(nil)
	if err != nil {
		return err
	}

	return s.createSessionDescription(ans)
}

func (s *receiveSession) onMessage() func(webrtc.DataChannelMessage) {
	return func(d webrtc.DataChannelMessage) {
		s.msgChannel <- d
	}
}

func (s *receiveSession) onClose() func() {
	return func() {
		close(s.sess.Done)
	}
}
