package main

import (
	"avtor.ru/bot/tg/internal/adapter"
	"avtor.ru/bot/tg/internal/port"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	envPath := "/Users/kirill/study/master/avtors_bot/src/bot/cmd/.env"

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	hostingAddress := os.Getenv("ADDRESS")
	if hostingAddress == "" {
		log.Fatal("ADDRESS environment variable is required")
	}

	portParam := os.Getenv("PORT")
	if portParam == "" {
		log.Fatal("ADDRESS environment variable is required")
	}

	serviceHost := fmt.Sprintf("http://%s:%s", hostingAddress, portParam)

	nspd, err := adapter.NewAnalyseServiceAdapter(serviceHost)
	if err != nil {
		log.Fatalf("Failed to create nspd client: %v", err)
	}

	um, err := adapter.NewUMServiceAdapter(serviceHost)
	if err != nil {
		log.Fatalf("Failed to create um client: %v", err)
	}

	bot, err := port.NewBot(botToken, nspd, um)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.Start(ctx)
}
