package main

import (
	"avtor.ru/bot/tg/internal/adapter"
	"avtor.ru/bot/tg/internal/port"
	"context"
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

	nspd, err := adapter.NewAnalyseServiceAdapter("http://127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Failed to create nspd client: %v", err)
	}

	bot, err := port.NewBot(botToken, nspd)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.Start(ctx)
}
