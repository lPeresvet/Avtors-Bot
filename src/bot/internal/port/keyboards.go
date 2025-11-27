package port

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackData string

var (
	AnalyseData    CallbackData = "analyse"
	LikedListData  CallbackData = "likedList"
	LikeData       CallbackData = "like"
	UnlikeData     CallbackData = "unlike"
	CreateUserData CallbackData = "createUser"
	AllUsersData   CallbackData = "allUsers"
	DeleteUserData CallbackData = "deleteUser"
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
	umRow = tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить пользователя", CreateUserData.String()),
		tgbotapi.NewInlineKeyboardButtonData("Список пользователей", AllUsersData.String()),
	)
	noAuth = tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вы не авторизованы", "/start"),
	)
)

var MainMenuNoAuthKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	noAuth,
)

var MainMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	analyseRow,
	likedListRow,
)

var MainMenuKeyboardForAdmin = tgbotapi.NewInlineKeyboardMarkup(
	analyseRow,
	likedListRow,
	umRow,
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
			tgbotapi.NewInlineKeyboardButtonData("Удалить из избранного ❌️", fmt.Sprintf("%s+%s", UnlikeData.String(), zoneID)),
		),
	)
}

func GetUserMenuKeyboard(userID string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить пользователя ❌️", fmt.Sprintf("%s+%s", DeleteUserData.String(), userID)),
		),
	)
}
