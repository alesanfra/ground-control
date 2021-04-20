package web

import (
	"fmt"
	"log"

	"github.com/alesanfra/ground-control/scanner"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Server rest api for the agent
type Server struct {
	devices scanner.DeviceMap
	port    uint
}

func NewWebServer(devices scanner.DeviceMap, port uint) *Server {
	return &Server{devices: devices, port: port}
}

// Start starts http server
func (s *Server) Start() {
	app := fiber.New()

	api := app.Group("/api", logger.New())
	api.Get("/status", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	v1 := api.Group("/v1")
	v1.Get("/devices", func(c *fiber.Ctx) error {
		return c.JSON(s.devices.AsList())
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", s.port)))
}
