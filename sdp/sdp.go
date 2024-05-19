package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pion/webrtc/v2"
)

const (
	name            = "sdp server"
	version         = "1.0.0"
	senderBuffSize  = 16384
	bufferThreshold = 512 * 1024
)

type app struct {
	name    string
	version string
	receive bool
	send    bool
	file    string
	output  string
	// secret  string
}

func NewApp() *app {
	return &app{
		name:    name,
		version: version,
		receive: false,
		send:    false,
		file:    "",
	}
}

type sendConfig struct {
	commConfiguration
	Stream io.Reader // The Stream to read from
}

type receiveConfig struct {
	commConfiguration
	Stream io.Writer // The Stream to write to
}

type commConfiguration struct {
	SDPProvider  io.Reader   // The SDP reader
	SDPOutput    io.Writer   // The SDP writer
	OnCompletion compHandler // Handler to call on session completion
	STUN         string      // Custom STUN server
}

type session struct {
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

type stats struct {
	lock      *sync.RWMutex
	nbBytes   uint64
	timeStart time.Time
	timeStop  time.Time

	timePause  time.Time
	timePaused time.Duration
}

type compHandler func()

func newStats() *stats {
	return &stats{
		lock: &sync.RWMutex{},
	}
}

func main() {
	var port string
	flag.StringVar(&port, "port", "9091", "listen port")

	a := NewApp()
	app := fiber.New()

	app.Use(cors.New())
	app.Get("/sdp", a.sdp)

	app.Listen(fmt.Sprintf(":%v", port))
}

func (a *app) sdp(c *fiber.Ctx) error {
	cf := sendConfig{
		Stream: nil,
		commConfiguration: commConfiguration{
			OnCompletion: func() {},
			STUN:         "stun.l.google.com:19302",
		},
	}

	s := newSession(cf)
	s.a = a

	err := s.start()
}

func newSession(conf sendConfig) *session {
	s := &session{
		sess: inSession{
			sdpInput:     conf.SDPProvider,
			sdpOutput:    conf.SDPOutput,
			Done:         make(chan struct{}),
			NetworkStats: newStats(),
			stunServers:  []string{fmt.Sprintf("stun:%s", c.STUN)},
			onCompletion: func() {},
		},
		stream:       conf.Stream,
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

func (s *session) start() error {
	return nil
}
