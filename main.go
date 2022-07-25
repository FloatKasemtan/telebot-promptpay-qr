package main

import (
	"MixkoPay/utils/config"
	"github.com/Frontware/promptpay"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

var InitKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Pay"),
	),
)

var SelectKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(KBank),
	), tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Scb),
	), tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(Bualuang),
	), tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Cancel"),
	),
)

func main() {
	bot, err := tgbotapi.NewBotAPI(config.C.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	isOpenPayment := false
	isSelectBank := false
	var promptPayBank string
	var amount float64

	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		if !isOpenPayment {
			switch strings.ToLower(update.Message.Text) {
			case "pay":
				msg.ReplyMarkup = SelectKeyboard
				isOpenPayment = true
			}
		} else {
			if !isSelectBank {
				switch strings.ToLower(update.Message.Text) {
				case "cancel":
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					isOpenPayment = false
				case "kbank":
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					promptPayBank = config.C.PromptPayKbankId
					isSelectBank = true
				case "scb":
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					promptPayBank = config.C.PromptPayKbankId
					isSelectBank = true
				case "bualuang":
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					promptPayBank = config.C.PromptPayKbankId
					isSelectBank = true
				}
				msg.Text = "Please enter your payment amount"
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				print(msg.Text)
				convertedAmount, err := strconv.Atoi(msg.Text)
				if err != nil {
					msg.Text = "Please enter a valid amount"
					msg.ReplyMarkup = SelectKeyboard
					break
				}
				amount = float64(convertedAmount)

				// * Generate QR code

				payment := promptpay.PromptPay{
					PromptPayID: promptPayBank, // Tax-ID/ID Card/E-Wallet
					Amount:      amount,        // Positive amount
				}

				// * Generate string to be use in QRCode
				qrcode, _ := payment.Gen()
				print(`---------------------------------------------------------` + qrcode + `-----------------------------------------------------`)
				// * Send QR code to user
				url := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL("https://chart.googleapis.com/chart?cht=qr&chs=500x500&chl=" + qrcode))

				mediaGroup := tgbotapi.NewMediaGroup(update.Message.Chat.ID, []interface{}{
					url,
				})

				if _, err := bot.Send(&mediaGroup); err != nil {
					log.Println("-------------------------------------------------------------------")
					log.Println(err)
					log.Println("-------------------------------------------------------------------")
				}
			}
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
