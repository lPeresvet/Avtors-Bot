package port

import (
	"avtor.ru/bot/client"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	api            *tgbotapi.BotAPI
	analyseService AnalyseService
}

type AnalyseService interface {
	Analyse(ctx context.Context, zoneID string) (*client.ZoneDetails, error)
}

func NewBot(token string, analyseService AnalyseService) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:            api,
		analyseService: analyseService,
	}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	b.api.Debug = false
	log.Printf("Authorized on account %s", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if err := b.handleMessage(ctx, update.Message); err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось обработать команду😔. Попробуйте позже")
				b.api.Send(msg)

				log.Printf("Error while handling message: %v", err)
			}
		}
	}

	return nil
}

func (b *Bot) handleMessage(ctx context.Context, message *tgbotapi.Message) error {
	switch message.Command() {
	case "start":
		b.sendWelcome(message.Chat.ID)
	case "analise":
		return b.analise(ctx, message.Chat.ID, message.CommandArguments())
	default:
		b.sendMainMenu(message.Chat.ID)
	}

	return nil
}

func (b *Bot) sendWelcome(chatID int64) {
	text := `👋 Добро пожаловать!

Я бот анализа земли. Пришлите /analise для анализа`

	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}

func (b *Bot) analise(ctx context.Context, chatID int64, zoneID string) error {
	//TODO validate zoneID
	log.Printf("Analyse zone with id %s", zoneID)
	zone, err := b.analyseService.Analyse(ctx, zoneID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID, FormatZone(zone))
	b.api.Send(msg)

	return nil
}

func (b *Bot) sendMainMenu(chatID int64) {
	text := "Главное меню:\n/analise - анализ земли"
	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}

func (b *Bot) Stop() {
	b.api.StopReceivingUpdates()
}
