package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pion/webrtc/v2"
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
			sdpInput:     c.SDPProvider,
			sdpOutput:    c.SDPOutput,
			Done:         make(chan struct{}),
			NetworkStats: newStats(),
			stunServers:  []string{fmt.Sprintf("stun:%s", c.STUN)},
			onCompletion: func() {},
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

func (s *inSession) createDataChannel(c *webrtc.DataChannelInit) (*webrtc.DataChannel, error) {
	return s.peerConnection.CreateDataChannel("data", c)
}

func (s *sendSession) close(call bool) {
	if !call {
		s.dataChannel.Close()
	}

	s.doneCheckLock.Lock()
	if s.doneCheck {
		s.doneCheckLock.Unlock()
		return
	}
	s.doneCheck = true
	s.doneCheckLock.Unlock()

	s.dumpStats()
	close(s.sess.Done)
}

// TODO: string implementation for below
func (s *sendSession) dumpStats() {
	fmt.Printf(`
Disk    : %v
Network : %v
`, s.readingStats, s.sess.NetworkStats)
}

func (s *sendSession) onBufferedAmountLow() func() {
	return func() {
		d := <-s.output
		if d.n != 0 {
			s.msgToBeSent = append(s.msgToBeSent, d)
		} else if len(s.msgToBeSent) == 0 && s.dataChannel.BufferedAmount() == 0 {
			s.sess.NetworkStats.stop()
			s.close(false)
			return
		}

		speed := s.sess.NetworkStats.bandwidth()
		fmt.Printf("%.2f MB/s\r", speed)

		for len(s.msgToBeSent) != 0 {
			cur := s.msgToBeSent[0]

			if err := s.dataChannel.Send(cur.buff); err != nil {
				logrus.Errorf("Error, cannot send to client: %v\n", err)
				return
			}

			s.sess.NetworkStats.addBytes(uint64(cur.n))
			s.msgToBeSent = s.msgToBeSent[1:]
		}
	}
}

func (s *sendSession) writeToNetwork() {
	s.dataChannel.OnBufferedAmountLow(s.onBufferedAmountLow())

	<-s.stopSending
	s.dataChannel.OnBufferedAmountLow(nil)
	s.sess.NetworkStats.pause()
	logrus.Infof("Pausing network I/O... (remaining at least %v packets)\n", len(s.output))
}

func (s *sendSession) onOpenHandler() func() {
	return func() {
		s.sess.NetworkStats.start()

		logrus.Info("start sending data ...")
		defer logrus.Info("stop sending data ...")

		s.writeToNetwork()
	}
}

func (s *sendSession) onCloseHandler() func() {
	return func() {
		s.close(true)
	}
}

func (s *sendSession) createDataChannel() error {
	o := true
	mxPacketLifeTime := uint16(10000)
	dChannel, err := s.sess.createDataChannel(&webrtc.DataChannelInit{
		Ordered:           &o,
		MaxPacketLifeTime: &mxPacketLifeTime,
	})

	if err != nil {
		return err
	}

	s.dataChannel = dChannel
	s.dataChannel.OnBufferedAmountLow(s.onBufferedAmountLow())
	s.dataChannel.SetBufferedAmountLowThreshold(bufferThreshold)
	s.dataChannel.OnOpen(s.onOpenHandler())
	s.dataChannel.OnClose(s.onCloseHandler())

	return nil
}

func (s *inSession) createOffer() error {
	ans, err := s.peerConnection.CreateOffer(nil)
	if err != nil {
		return err
	}

	return s.createSessionDescription(ans)
}

func (s *inSession) createSessionDescription(d webrtc.SessionDescription) error {
	if err := s.peerConnection.SetLocalDescription(d); err != nil {
		return err
	}

	d.SDP = stripSDP(d.SDP)

	res, err := encode(d)
	if err != nil {
		return err
	}

	fmt.Fprintf(s.sdpOutput, "%s\n", res)
	return nil
}

func (s *sendSession) initialize() error {
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

func (s *sendSession) readFile() {
	logrus.Info("start reading data ...")
	s.readingStats.start()

	defer func() {
		s.readingStats.pause()
		logrus.Info("reading data paused ...")
		close(s.output)
	}()

	for {
		s.dataBuff = s.dataBuff[:cap(s.dataBuff)]
		n, err := s.stream.Read(s.dataBuff)
		if err != nil {
			switch {
			case err == io.EOF:
				s.readingStats.stop()
				logrus.Debugf("Got EOF after %v bytes!\n", s.readingStats.bytes())
			default:
				logrus.Errorf("reading error: %v", err)
			}
			return
		}

		s.dataBuff = s.dataBuff[:n]
		s.readingStats.addBytes(uint64(n))

		s.output <- outputMsg{
			n:    n,
			buff: append([]byte(nil), s.dataBuff...),
		}
	}
}

func (s *inSession) readSDP() error {
	var sdp webrtc.SessionDescription

	for {
		enc, err := readStream(s.sdpInput)
		if err == nil {
			if err := decode(enc, &sdp); err == nil {
				break
			}
		}
		logrus.Println("Invalid SDP, try again ...")
	}

	return s.peerConnection.SetRemoteDescription(sdp)
}

// TODO
func (s *sendSession) start() error {
	if err := s.initialize(); err != nil {
		return err
	}

	go s.readFile()
	if err := s.sess.readSDP(); err != nil {
		logrus.Errorln(err)
		return err
	}

	<-s.sess.Done
	s.sess.onCompletion()
	return nil
}
