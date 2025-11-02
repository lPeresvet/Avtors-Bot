package port

import (
	"avtor.ru/bot/client"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackData string

var (
	AnalyseData   CallbackData = "analyse"
	LikedListData CallbackData = "likedList"
	LikeData      CallbackData = "like"
)

func (c *CallbackData) String() string {
	return string(*c)
}

var ZoneMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить в избранное ❤️", LikeData.String()),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Провести анализ 🔍", AnalyseData.String()),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Cписок избранных участков 📋", LikedListData.String()),
	),
)

var MainMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Провести анализ 🔍", AnalyseData.String()),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Cписок избранных участков 📋", LikedListData.String()),
	),
)

func FormatZone(zone *client.ZoneDetails) string {
	return fmt.Sprintf("Кадастровый номер: %v\nФормат собственности: %v\nВид использования: %v", zone.Id, zone.PropertyType, zone.PermittedUsage)
}
