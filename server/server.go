package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type App struct {
	LocalDSP  string
	RemoteDSP string
	FilePath  string
}

type SendRequest struct {
	Remote string `json:"remote"`
	Local  string `json:"local"`
}

type SendResponse struct{}

type ReceiveRequest struct {
	Remote string `json:"remote"`
	Local  string `json:"local"`
}

type ReceiveResponse struct{}

type DSPResponse struct {
	DSP string `json:"dsp"`
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	ap := App{}

	app.Post("/send", ap.sendHandler) // local -> remote

	app.Post("/receive", ap.receiveHandler) // remote -> local

	app.Get("/dsp", ap.getDSP)

	app.Post("/dsp", ap.postDSP)

	app.Listen(":3000")
}

func (app *App) sendHandler(c *fiber.Ctx) error { return nil }

func (app *App) receiveHandler(c *fiber.Ctx) error { return nil }

func (app *App) getDSP(c *fiber.Ctx) error {
	return nil
}

func (app *App) postDSP(c *fiber.Ctx) error { return nil }
