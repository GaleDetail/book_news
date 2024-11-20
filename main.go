package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	const (
		url      = "https://chytomo.com/?s=%D0%BA%D0%BD%D0%B8%D0%B3%D0%B0"
		botToken = "your_code"
	)

	// Ініціалізація Telegram-бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Налаштування оновлень
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // Повідомлення
			if update.Message.Text == "/start" {
				// Відправка повідомлення з кнопкою
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Натисніть кнопку, щоб побачити останні статті:")
				msg.ReplyMarkup = createKeyboard()
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Error sending start message: %v", err)
				}
			}
		} else if update.CallbackQuery != nil { // Натискання кнопки
			callback := update.CallbackQuery
			chatID := callback.Message.Chat.ID // ID чату користувача

			if callback.Data == "show_articles" {
				// Отримання статей
				articles, err := FetchLatestArticles(url)
				if err != nil {
					log.Printf("Error fetching articles: %v", err)
					// Відповідь на натискання кнопки
					bot.Send(tgbotapi.NewMessage(chatID, "Помилка отримання статей"))
					continue
				}

				merged := mergeJson(articles)
				fmt.Println(&merged, "<- ось і все")
				// Відправка тільки нових статей
				for _, article := range articles {
					// Перевірка: чи є стаття в уже збережених (merged)
					if IsArticleInMerged(article, merged) {
						log.Printf("Article already sent: %s\n", article.Title)
						continue
					}

					// Формування посту
					post := CreatePost(article)

					// Надсилання посту до Telegram
					if err := SendToTelegram(bot, chatID, post, article.Image); err != nil {
						log.Fatalf("Error sending post to Telegram: %v", err)
					}

					log.Println("Post successfully sent to Telegram!")

				}

				// Відповідь на натискання кнопки
				//photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL("https://avatars.githubusercontent.com/u/165570728?v=4&size=40"))
				msg := tgbotapi.NewMessage(chatID, "Статті надіслані!")
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Error sending confirmation: %v", err)
				}
			}
		}
	}
}

func createKeyboard() tgbotapi.InlineKeyboardMarkup {
	button := tgbotapi.NewInlineKeyboardButtonData("Показати статті", "show_articles")
	row := []tgbotapi.InlineKeyboardButton{button}
	return tgbotapi.NewInlineKeyboardMarkup(row)
}

func mergeJson(newArticles []*Article) []*StoredArticle {
	existingArticles, err := LoadArticles("articles.json")
	if err != nil {
		log.Fatalf("Error loading articles: %v", err)
	}

	mergedArticles, err := MergeArticles(existingArticles, newArticles)
	if err != nil {
		log.Fatalf("Error merging articles: %v", err)
	}

	err = SaveArticles("articles.json", mergedArticles)
	if err != nil {
		log.Fatalf("Error saving articles: %v", err)
	}
	return mergedArticles
}
func IsArticleInMerged(article *Article, merged []*StoredArticle) bool {
	for _, stored := range merged {
		if stored.URL == article.Link {
			return true // Стаття вже є у списку
		}
	}
	return false // Стаття нова
}
