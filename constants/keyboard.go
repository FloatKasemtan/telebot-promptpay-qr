package constants

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var SelectKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(KBank),
	), tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Scb),
	), tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Bualuang),
	), tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Cancel),
	),
)
