package port

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackData string

var (
	AnalyseData   CallbackData = "analyse"
	LikedListData CallbackData = "likedList"
	LikeData      CallbackData = "like"
	UnikeData     CallbackData = "unlike"
)

func (c *CallbackData) String() string {
	return string(*c)
}

var (
	analyseRow = tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Провести анализ 🔍", AnalyseData.String()),
	)
	likedListRow = tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Cписок избранных участков 📋", LikedListData.String()),
	)
)

var MainMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	analyseRow,
	likedListRow,
)

func GetZoneMenuKeyboard(zoneID string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить в избранное ❤️", fmt.Sprintf("%s+%s", LikeData.String(), zoneID)),
		),
		analyseRow,
		likedListRow,
	)
}

func GetLikedZoneMenuKeyboard(zoneID string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить из избранного ❌️", fmt.Sprintf("%s+%s", UnikeData.String(), zoneID)),
		),
	)
}
