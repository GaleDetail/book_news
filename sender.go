package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendToTelegram(bot *tgbotapi.BotAPI, chatID int64, post, imageURL string) error {
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(imageURL))
	photo.Caption = post
	photo.ParseMode = "Markdown"
	_, err := bot.Send(photo)
	if err != nil {
		return fmt.Errorf("failed to send post to Telegram: %w", err)
	}
	return nil
}
