package di

import (
	"avtor.ru/bot/analyse_service/internal/client"
	"avtor.ru/bot/analyse_service/internal/handlers"
	"avtor.ru/bot/server"
	"context"
	"github.com/labstack/echo/v4"
)

type Container struct {
	ctx  context.Context
	Echo *echo.Echo

	NSPDClient handlers.NSPDClient
	Server     server.ServerInterface
}

func NewContainer() *Container {
	return &Container{
		ctx: context.Background(),
	}
}

func (c *Container) GetEcho() *echo.Echo {
	if c.Echo == nil {
		c.Echo = echo.New()
	}

	return c.Echo
}

func (c *Container) GetNSPDClient() handlers.NSPDClient {
	if c.NSPDClient == nil {
		c.NSPDClient = client.NewNSDPClient()
	}

	return c.NSPDClient
}

func (c *Container) GetService() server.ServerInterface {
	if c.Server == nil {
		c.Server = handlers.NewAnalyseService(c.ctx, c.GetNSPDClient())
	}

	return c.Server
}
