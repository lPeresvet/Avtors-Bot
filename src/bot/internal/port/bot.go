package port

import (
	"avtor.ru/bot/client"
	"avtor.ru/bot/tg/internal/adapter"
	"avtor.ru/bot/tg/internal/usecase/zones"
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"sync"
)

type State string

var (
	AnalyseState    State = "analyse"
	CreateUserState State = "create_user"
	NoState         State = "no"
)

type Role string

var (
	AdminRole     Role = "admin"
	AnalyserRole  Role = "analyser"
	UndefinedRole Role = "undefined"
)

type Bot struct {
	api            *tgbotapi.BotAPI
	analyseService AnalyseService
	userService    UserService

	mu     sync.Mutex
	states map[int64]State

	roles map[int64]Role
}

type AnalyseService interface {
	Analyse(ctx context.Context, zoneID string) (*client.ZoneDetails, error)
	GetLikes(ctx context.Context) (*client.Zones, error)
	LikeZone(ctx context.Context, userID int64, zoneID string) error
	UnlikeZone(ctx context.Context, userID int64, zoneID string) error
}

type UserService interface {
	CreateUser(ctx context.Context, username, role string) error
	GetUserRole(ctx context.Context, username string) (string, error)
	GetUsers(ctx context.Context) (*client.Users, error)
	DeleteUser(ctx context.Context, username string) error
}

func NewBot(token string, analyseService AnalyseService, umService UserService) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:            api,
		analyseService: analyseService,
		userService:    umService,
		states:         make(map[int64]State),
		roles:          make(map[int64]Role),
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
			data := strings.Split(update.CallbackQuery.Data, "+")
			cmd := data[0]
			payload := ""

			if len(data) > 1 {
				payload = data[1]
			}

			if err := b.handleQuery(ctx, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.From.UserName, cmd, payload); err != nil {
				log.Printf("Error while handling callback: %v", err)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось обработать команду😔. Попробуйте позже")
				b.api.Send(msg)

				continue
			}

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Done!")
			if _, err := b.api.Request(callback); err != nil {
				log.Printf("Error while processing callback: %v", err)
			}
		}
	}

	return nil
}

func (b *Bot) handleMessage(ctx context.Context, message *tgbotapi.Message) error {
	if err := b.ProceedUserRole(ctx, message.Chat.ID, message.From.UserName); err != nil {
		return fmt.Errorf("failed to proceed user role: %w", err)
	}

	switch b.getUserState(message.Chat.ID) {
	case AnalyseState:
		b.clearUserState(message.Chat.ID)
		return b.analise(ctx, message.Chat.ID, message.Text)
	case CreateUserState:
		b.clearUserState(message.Chat.ID)

		return b.createUser(ctx, message.Chat.ID, message.Text)
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

func (b *Bot) handleQuery(ctx context.Context, chatID int64, username string, callbackCMD, callbackPayload string) error {
	if err := b.ProceedUserRole(ctx, chatID, username); err != nil {
		return fmt.Errorf("failed to proceed user role: %w", err)
	}

	msgs := make([]tgbotapi.MessageConfig, 0)

	switch CallbackData(callbackCMD) {
	case AnalyseData:
		outText := "Введите номер кадастрового участка"

		b.setUserState(chatID, AnalyseState)
		msgs = append(msgs, tgbotapi.NewMessage(chatID, outText))
	case CreateUserData:
		outText := "Введите пользователя в формате: tg_имя_пользователя - роль\n Возможны роли Admin или Analyser"

		b.setUserState(chatID, CreateUserState)
		msgs = append(msgs, tgbotapi.NewMessage(chatID, outText))
	case AllUsersData:
		users, err := b.userService.GetUsers(ctx)
		if err != nil {
			return fmt.Errorf("failed to get likes: %v", err)
		}

		for _, user := range *users {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Имя пользователя: @%s\nРоль: %s", *user.Username, *user.Role))
			msg.ReplyMarkup = GetUserMenuKeyboard(*user.Username)

			msgs = append(msgs, msg)
		}

		msg := tgbotapi.NewMessage(chatID, "Меню:")

		msg.ReplyMarkup = b.GetMainMenuBasedOnRole(chatID)
		msgs = append(msgs, msg)
	case DeleteUserData:
		outText := "Пользователь успешно удален ✅"

		if err := b.userService.DeleteUser(ctx, callbackPayload); err != nil {
			outText = "Не удалось удалить пользователя😔"

			log.Printf("failed to delete user: %v", err)
		}

		msg := tgbotapi.NewMessage(chatID, outText)
		msg.ReplyMarkup = b.GetMainMenuBasedOnRole(chatID)

		msgs = append(msgs, msg)

	case LikedListData:
		likes, err := b.analyseService.GetLikes(ctx)
		if err != nil {
			return fmt.Errorf("failed to get likes: %v", err)
		}

		for _, like := range *likes {
			msg := tgbotapi.NewMessage(chatID, like.Id)
			msg.ReplyMarkup = GetLikedZoneMenuKeyboard(like.Id)

			msgs = append(msgs, msg)
		}

		msg := tgbotapi.NewMessage(chatID, "Меню:")

		msg.ReplyMarkup = b.GetMainMenuBasedOnRole(chatID)
		msgs = append(msgs, msg)

	case LikeData:
		outText := "Участок добавлен в избранное ✅"
		if err := b.analyseService.LikeZone(ctx, chatID, callbackPayload); err != nil {
			if errors.Is(err, adapter.ErrorLikeZone) {
				outText = "Не удалось добавить участок в избранное 😔"
			}

			log.Printf("failed to like zone: %v", err)
		}

		msg := tgbotapi.NewMessage(chatID, outText)

		msg.ReplyMarkup = b.GetMainMenuBasedOnRole(chatID)

		msgs = append(msgs, msg)
	case UnlikeData:
		outText := "Участок удален из избранного ✅"

		if err := b.analyseService.UnlikeZone(ctx, chatID, callbackPayload); err != nil {
			outText = "Не удалось удалить участок из избранного 😔"

			log.Printf("failed to unlike zone: %v", err)
		}

		msg := tgbotapi.NewMessage(chatID, outText)
		msg.ReplyMarkup = b.GetMainMenuBasedOnRole(chatID)

		msgs = append(msgs, msg)
	}

	for _, msg := range msgs {
		if _, err := b.api.Send(msg); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) sendWelcome(chatID int64) {
	text := `👋 Добро пожаловать!

Я бот анализа земли. Вот что я могу ⬇️`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = b.GetMainMenuBasedOnRole(chatID)
	b.api.Send(msg)
}

func (b *Bot) analise(ctx context.Context, chatID int64, zoneID string) error {
	var msg tgbotapi.MessageConfig
	log.Printf("Analyse zone with id %s", zoneID)

	if !zones.ValidateZone(zoneID) {
		msg = tgbotapi.NewMessage(chatID, "Кадастровый номер участка невалиден ⚠️")
		msg.ReplyMarkup = b.GetMainMenuBasedOnRole(chatID)
	} else {
		zone, err := b.analyseService.Analyse(ctx, zoneID)
		if err != nil {
			return err
		}

		msg = tgbotapi.NewMessage(chatID, FormatZone(zone))
		msg.ReplyMarkup = GetZoneMenuKeyboard(zoneID)
		msg.ParseMode = tgbotapi.ModeMarkdown
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
	msg.ReplyMarkup = b.GetMainMenuBasedOnRole(chatID)

	b.api.Send(msg)
}

func (b *Bot) Stop() {
	b.api.StopReceivingUpdates()
}

func (b *Bot) ProceedUserRole(ctx context.Context, chatID int64, username string) error {
	log.Printf("ProceedUserRole: username=%s, chatID=%d", username, chatID)

	role, err := b.userService.GetUserRole(ctx, username)
	if err != nil {
		b.clearUserRole(chatID)
		log.Printf("failed to get user role: %v\n", err)

		//return fmt.Errorf("failed to retrieve user role: %v", err)
		return nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	b.roles[chatID] = Role(role)

	return nil
}

func (b *Bot) clearUserRole(chatID int64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.roles, chatID)
}

func (b *Bot) GetMainMenuBasedOnRole(chatID int64) tgbotapi.InlineKeyboardMarkup {
	b.mu.Lock()
	defer b.mu.Unlock()

	role := Role(strings.ToLower(string(b.roles[chatID])))
	if role == AdminRole {
		return MainMenuKeyboardForAdmin
	} else if role == AnalyserRole {
		return MainMenuKeyboard
	}

	return MainMenuNoAuthKeyboard
}

func (b *Bot) createUser(ctx context.Context, chatID int64, msgData string) error {
	var text string
	data := strings.Split(msgData, " - ")
	text = "Пользователь создан"

	if len(data) == 2 {
		newRole := strings.ToLower(data[1])
		log.Printf(newRole)
		if !(newRole == string(AdminRole) || newRole == string(AnalyserRole)) {
			text = "Невалидная роль"
		} else if err := b.userService.CreateUser(ctx, data[0], data[1]); err != nil {
			text = "Не получилось создать нового пользователя"

			log.Printf("Failed to create user: %v\n", err)
		}
	} else {
		text = "Невалидный формат запроса"
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = b.GetMainMenuBasedOnRole(chatID)

	b.api.Send(msg)

	return nil
}
