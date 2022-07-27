package main

import (
	"MixkoPay/constants"
	"MixkoPay/utils/config"
	"github.com/Frontware/promptpay"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

// Initialize global variables
var isOpenPayment = false
var isSelectBank = false

func main() {
	// Initialize telegram bot
	bot, err := tgbotapi.NewBotAPI(config.C.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Receive updates
	updates := bot.GetUpdatesChan(u)

	var promptPayBank string
	var amount float64

	for update := range updates {
		// ignore non-Message updates
		if update.Message == nil {
			continue
		}

		// Initialize new message instance
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		if !isOpenPayment {
			if strings.ToLower(update.Message.Text) == "pay" {
				msg.ReplyMarkup = constants.SelectKeyboard
				isOpenPayment = true
			}
		} else {
			if !isSelectBank {
				promptPayBank = SelectBank(update, &msg, promptPayBank)

				if isSelectBank {
					msg.Text = "Please enter your payment amount"
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				}
			} else {
				// Convert amount to number
				if amount, isError := ConvertAmount(&msg, amount); !isError {
					// * Generate QR code
					qrcode := GenerateQRPayment(promptPayBank, amount)

					// * Send QR code to user
					SendQRPayment(qrcode, update, bot)
				}

			}
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func SendQRPayment(qrcode string, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	url := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL("https://chart.googleapis.com/chart?cht=qr&chs=500x500&chl=" + qrcode))

	mediaGroup := tgbotapi.NewMediaGroup(update.Message.Chat.ID, []interface{}{
		url,
	})

	if _, err := bot.Request(&mediaGroup); err != nil {
		log.Panic(err)
	}
}

func GenerateQRPayment(promptPayBank string, amount float64) string {
	payment := promptpay.PromptPay{
		PromptPayID: promptPayBank, // Tax-ID/ID Card/E-Wallet
		Amount:      amount,        // Positive amount
	}
	isSelectBank = false
	isOpenPayment = false
	qrcode, _ := payment.Gen()
	return qrcode
}

func ConvertAmount(msg *tgbotapi.MessageConfig, amount float64) (float64, bool) {
	convertedAmount, err := strconv.Atoi(msg.Text)
	if err != nil {
		msg.Text = "Please enter a valid amount"
		return 0, true
	}
	amount = float64(convertedAmount)
	return amount, false
}

func SelectBank(update tgbotapi.Update, msg *tgbotapi.MessageConfig, promptPayBank string) string {
	switch strings.ToLower(update.Message.Text) {
	case constants.Cancel:
		isOpenPayment = false
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	case constants.KBank:
		promptPayBank = config.C.PromptPayKbankId
		isSelectBank = true
	case constants.Scb:
		msg.Text = "Not currently support SCB account"
		//promptPayBank = config.C.PromptPayKbankId
		//isSelectBank = true
	case constants.Bualuang:
		msg.Text = "Not currently support Bualuang account"
		//promptPayBank = config.C.PromptPayKbankId
		//isSelectBank = true
	}
	return promptPayBank
}
