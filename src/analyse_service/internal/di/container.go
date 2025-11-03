package di

import (
	"avtor.ru/bot/analyse_service/internal/client"
	"avtor.ru/bot/analyse_service/internal/handlers"
	"avtor.ru/bot/analyse_service/internal/model"
	"avtor.ru/bot/analyse_service/internal/repository"
	"avtor.ru/bot/server"
	"context"
	"github.com/labstack/echo/v4"
)

type Container struct {
	ctx  context.Context
	Echo *echo.Echo

	NSPDClient handlers.NSPDClient
	Server     server.ServerInterface
	Repository *repository.Repository
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

func (c *Container) GetService() (server.ServerInterface, error) {
	if c.Server == nil {
		repo, err := c.GetRepository()
		if err != nil {
			return nil, err
		}

		c.Server = handlers.NewAnalyseService(c.ctx, c.GetNSPDClient(), repo)
	}

	return c.Server, nil
}

func (c *Container) GetRepository() (*repository.Repository, error) {
	config := model.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "developer",
		Password: "developer",
		DBName:   "liked_zones_service",
	}

	if c.Repository == nil {
		repo, err := repository.NewRepository(&config)
		if err != nil {
			return nil, err
		}

		c.Repository = repo
	}

	return c.Repository, nil
}
