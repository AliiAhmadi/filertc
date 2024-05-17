package main

import (
	"github.com/gofiber/fiber/v2"
)

type App struct {
	localDSP  string
	remoteDSP string
}

type SendRequest struct{}

type SendResponse struct{}

type ReceiveRequest struct{}

type ReceiveResponse struct{}

func main() {
	app := fiber.New()
	ap := App{}

	app.Post("/send", ap.sendHandler) // local -> remote

	app.Post("/receive", ap.receiveHandler) // remote -> local

	app.Get("/dsp", ap.getDSP)

	app.Post("/dsp", ap.postDSP)

	app.Listen(":3000")
}

func (app *App) sendHandler(c *fiber.Ctx) error { return nil }

func (app *App) receiveHandler(c *fiber.Ctx) error { return nil }

func (app *App) getDSP(c *fiber.Ctx) error { return nil }

func (app *App) postDSP(c *fiber.Ctx) error { return nil }
