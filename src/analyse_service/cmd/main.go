package main

import (
	"log"

	_ "github.com/labstack/echo/v4"

	"avtor.ru/bot/analyse_service/internal/di"
	"avtor.ru/bot/server"
)

func main() {
	container := di.NewContainer()
	echoServer := container.GetEcho()

	server.RegisterHandlers(echoServer, container.GetService())

	log.Fatal(echoServer.Start("127.0.0.1:8080"))
}
