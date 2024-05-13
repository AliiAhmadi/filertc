package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pion/webrtc/v4"
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

func (s *sendSession) start() error {
	return nil
}
