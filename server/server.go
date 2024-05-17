package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	app.Get("/sdp", getSdp)

	app.Listen(":3000")
}

func send(c *fiber.Ctx) error {
	var req Request
	return nil
}

func receive(c *fiber.Ctx) error {
	var req Request
	return nil
}

func getSdp(c *fiber.Ctx) error {
	return nil
}
