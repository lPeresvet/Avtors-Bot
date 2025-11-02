package port

import (
	"avtor.ru/bot/client"
	"avtor.ru/bot/tg/internal/usecase/zones"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
)

type State string

var (
	AnalyseState State = "analyse"
	NoState      State = "no"
)

type Bot struct {
	api            *tgbotapi.BotAPI
	analyseService AnalyseService

	mu     sync.Mutex
	states map[int64]State
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
		states:         make(map[int64]State),
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
		} else if update.CallbackQuery != nil {
			if err := b.handleQuery(ctx, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data); err != nil {
				log.Printf("Error while handling callback: %v", err)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось обработать команду😔. Попробуйте позже")
				b.api.Send(msg)

				continue
			}

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := b.api.Request(callback); err != nil {
				log.Printf("Error while processing callback: %v", err)
			}
		}
	}

	return nil
}

func (b *Bot) handleMessage(ctx context.Context, message *tgbotapi.Message) error {
	switch b.getUserState(message.Chat.ID) {
	case AnalyseState:
		b.clearUserState(message.Chat.ID)
		return b.analise(ctx, message.Chat.ID, message.Text)
	}

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

func (b *Bot) handleQuery(_ context.Context, chatID int64, callbackData string) error {
	var (
		msg tgbotapi.MessageConfig
	)

	switch CallbackData(callbackData) {
	case AnalyseData:
		outText := "Введите номер кадастрового участка"

		b.setUserState(chatID, AnalyseState)
		msg = tgbotapi.NewMessage(chatID, outText)
	case LikedListData:
		outText := "Пока не реализовано. Тут будет список помеченных участков"
		msg = tgbotapi.NewMessage(chatID, outText)
	case LikeData:
		outText := "Участок добавлен в избранное ✅"
		msg = tgbotapi.NewMessage(chatID, outText)

		//TODO handle like logic

		msg.ReplyMarkup = MainMenuKeyboard
	}

	if _, err := b.api.Send(msg); err != nil {
		return err
	}

	return nil
}

func (b *Bot) sendWelcome(chatID int64) {
	text := `👋 Добро пожаловать!

Я бот анализа земли. Вот что я могу ⬇️`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = MainMenuKeyboard
	b.api.Send(msg)
}

func (b *Bot) analise(ctx context.Context, chatID int64, zoneID string) error {
	var msg tgbotapi.MessageConfig
	log.Printf("Analyse zone with id %s", zoneID)

	if !zones.ValidateZone(zoneID) {
		msg = tgbotapi.NewMessage(chatID, "Кадастровый номер участка невалиден ⚠️")
		msg.ReplyMarkup = MainMenuKeyboard
	} else {
		zone, err := b.analyseService.Analyse(ctx, zoneID)
		if err != nil {
			return err
		}

		msg = tgbotapi.NewMessage(chatID, FormatZone(zone))
		msg.ReplyMarkup = ZoneMenuKeyboard
	}

	b.api.Send(msg)

	return nil
}

func (b *Bot) setUserState(userID int64, state State) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.states[userID] = state
}

func (b *Bot) getUserState(userID int64) State {
	b.mu.Lock()
	defer b.mu.Unlock()

	if state, ok := b.states[userID]; ok {
		return state
	}

	return NoState
}

func (b *Bot) clearUserState(userID int64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.states, userID)
}

func (b *Bot) sendMainMenu(chatID int64) {
	text := "Главное меню"
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = MainMenuKeyboard

	b.api.Send(msg)
}

func (b *Bot) Stop() {
	b.api.StopReceivingUpdates()
}
