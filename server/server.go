package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const (
	senderBuffSize  = 16384
	bufferThreshold = 512 * 1024
	stun            = "stun.l.google.com:19302"
)

type Request struct {
	Name    string `json:"name"`
	Receive bool   `json:"receive"`
	Send    bool   `json:"send"`
	File    string `json:"file"`
	Output  string `json:"output"`
}

// OK    => 1
// Error => 2
type Response struct {
	Status int   `json:"status"`
	Size   int64 `json:"size"`
}

type SDPResponse struct {
	SDP string `json:"sdp"`
}

func main() {
	app := fiber.New()
	app.Use(cors.New())

	app.Post("/send", send)
	app.Post("/receive", receive)

	app.Listen(":3000")
}

func send(c *fiber.Ctx) error {
	var req Request

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	out, err := exec.Command("frtc", "--send", "--file", req.File).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The date is %s\n", out)

	return nil
}

func receive(c *fiber.Ctx) error {
	var req Request

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	return nil
}
