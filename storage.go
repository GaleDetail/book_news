package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type StoredArticle struct {
	URL  string    `json:"url"`
	Date time.Time `json:"date"`
}

// LoadArticles завантажує запощені статті
func LoadArticles(filename string) ([]*StoredArticle, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []*StoredArticle{}, nil // Повернути порожній список, якщо файл не існує
		}
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var articles []*StoredArticle
	err = json.NewDecoder(file).Decode(&articles)
	return articles, err
}

// SaveArticles зберігає запощені статті
func SaveArticles(filename string, articles []*StoredArticle) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(articles)
}

// FilterRecentArticles залишає тільки статті за останній тиждень
func FilterRecentArticles(articles []StoredArticle) []StoredArticle {
	var recent []StoredArticle
	weekAgo := time.Now().AddDate(0, 0, -7)

	for _, article := range articles {
		if article.Date.After(weekAgo) {
			recent = append(recent, article)
		}
	}
	return recent
}

// MergeArticles повертає статті, які є одночасно у збережених і нових
func MergeArticles(existing []*StoredArticle, newArticles []*Article) ([]*StoredArticle, error) {
	// Конвертуємо нові статті у `StoredArticle`
	storedNewArticles, err := ConvertArticlesToStored(newArticles)
	if err != nil {
		return nil, err
	}

	// Мапа для швидкого пошуку нових статей
	newArticleMap := make(map[string]struct{})
	for _, article := range storedNewArticles {
		newArticleMap[article.URL] = struct{}{}
	}
	fmt.Println(newArticleMap, "<-нові статті")
	fmt.Println(existing, "<-старі статті")
	// Залишаємо тільки ті статті, які є в обох списках
	var commonArticles []*StoredArticle
	for _, article := range existing {
		if _, exists := newArticleMap[article.URL]; exists {
			commonArticles = append(commonArticles, article)
		}
	}

	return commonArticles, nil
}

func ConvertArticlesToStored(articles []*Article) ([]*StoredArticle, error) {
	var storedArticles []*StoredArticle
	for _, article := range articles {

		// Парсимо дату у форматі "DD.MM.YYYY"
		parsedDate, err := time.Parse("02.01.2006", article.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date for article %s: %w", article.Link, err)
		}

		// Додаємо статтю до результату
		storedArticles = append(storedArticles, &StoredArticle{
			URL:  article.Link,
			Date: parsedDate,
		})
	}
	return storedArticles, nil
}
