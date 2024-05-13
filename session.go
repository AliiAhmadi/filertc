package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pion/webrtc/v4"
	"github.com/sirupsen/logrus"
)

type sendSession struct {
	sess          inSession
	stream        io.Reader
	initialized   bool
	dataChannel   *webrtc.DataChannel
	dataBuff      []byte
	msgToBeSent   []outputMsg
	stopSending   chan struct{}
	output        chan outputMsg
	doneCheckLock sync.Mutex
	doneCheck     bool
	readingStats  *stats
	a             *app
}

type inSession struct {
	Done           chan struct{}
	NetworkStats   *stats
	sdpInput       io.Reader
	sdpOutput      io.Writer
	peerConnection *webrtc.PeerConnection
	onCompletion   compHandler
	stunServers    []string
}

type outputMsg struct {
	n    int
	buff []byte
}

func newSendSession(c sendConfig) *sendSession {
	s := &sendSession{
		sess: inSession{
			sdpInput:  c.SDPProvider,
			sdpOutput: c.SDPOutput,
			Done:      make(chan struct{}),
			NetworkStats: &stats{
				lock: &sync.RWMutex{},
			},
			stunServers: []string{fmt.Sprintf("stun:%s", c.STUN)},
		},
		stream:       c.Stream,
		initialized:  false,
		dataBuff:     make([]byte, senderBuffSize),
		stopSending:  make(chan struct{}, 1),
		output:       make(chan outputMsg, senderBuffSize*10),
		doneCheck:    false,
		readingStats: newStats(),
	}

	if s.sess.sdpInput == nil {
		s.sess.sdpInput = os.Stdin
	}
	if s.sess.sdpOutput == nil {
		s.sess.sdpOutput = os.Stdout
	}

	return s
}

func (s *sendSession) onConnectionStateChange() func(webrtc.ICEConnectionState) {
	return func(connState webrtc.ICEConnectionState) {
		logrus.Infof("ICE Connection State has changed: %s\n", connState.String())
		if connState == webrtc.ICEConnectionStateDisconnected {
			s.stopSending <- struct{}{}
		}
	}
}

func (s *inSession) createConnection(changeFunction func(webrtc.ICEConnectionState)) error {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: s.stunServers,
			},
		},
	}

	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return err
	}

	s.peerConnection = pc
	pc.OnICEConnectionStateChange(changeFunction)
	return nil
}

func (s *sendSession) Initialize() error {
	if s.initialized {
		return nil
	}

	if err := s.sess.createConnection(s.onConnectionStateChange()); err != nil {
		logrus.Errorln(err)
		return err
	}

	if err := s.createDataChannel(); err != nil {
		logrus.Errorln(err)
		return err
	}

	if err := s.sess.createOffer(); err != nil {
		logrus.Errorln(err)
		return err
	}

	s.initialized = true
	return nil
}

// TODO
func (s *sendSession) start() error {
	if err := s.Initialize(); err != nil {
		return err
	}

	go s.readFile()
	if err := s.sess.ReadSDP(); err != nil {
		logrus.Errorln(err)
		return err
	}

	<-s.sess.Done
	s.sess.onCompletion()
	return nil
}
