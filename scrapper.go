package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
)

// Article представляє структуру для збереження інформації про статтю
type Article struct {
	Title   string
	Link    string
	Content string
	Image   string
	Date    string
}

// FetchHTML отримує HTML-документ за вказаним URL
func FetchHTML(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return doc, nil
}

// ParseArticle парсить HTML і повертає зріз Article
func ParseArticle(doc *goquery.Document) ([]*Article, error) {
	articles := doc.Find(".search-item")
	if articles.Length() == 0 {
		return nil, errors.New("no articles found")
	}

	var structArticles []*Article
	var parseErrors []error

	articles.Each(func(i int, selection *goquery.Selection) {
		article, err := extractArticle(selection, i)
		if err != nil {
			parseErrors = append(parseErrors, err)
			return
		}
		structArticles = append(structArticles, article)
	})

	if len(parseErrors) > 0 {
		return nil, fmt.Errorf("encountered %d parsing errors: %v", len(parseErrors), parseErrors)
	}

	return structArticles, nil
}

// extractArticle витягує дані статті з поточного блоку HTML
type ArticleError struct {
	Index int
	Field string
	Err   error
}

func (e *ArticleError) Error() string {
	return fmt.Sprintf("error at article index %d, field '%s': %v", e.Index, e.Field, e.Err)
}

// Регулярний вираз для знаходження дати у форматі "15.11.2024"
var dateRegex = regexp.MustCompile(`\d{2}\.\d{2}\.\d{4}`)

func extractArticle(selection *goquery.Selection, index int) (*Article, error) {
	title := selection.Find(".caption").Text()
	if title == "" {
		return &Article{}, &ArticleError{Index: index, Field: "title", Err: errors.New("missing value")}
	}

	link, exists := selection.Find("a").Attr("href")
	if !exists {
		return &Article{}, &ArticleError{Index: index, Field: "link", Err: errors.New("missing value")}
	}

	content := selection.Find(".caption--black").Text()
	if content == "" {
		return &Article{}, &ArticleError{Index: index, Field: "content", Err: errors.New("missing value")}
	}

	imgSrc, exists := selection.Find("img").Attr("src")
	if !exists {
		return &Article{}, &ArticleError{Index: index, Field: "image", Err: errors.New("missing value")}
	}

	rawDate := selection.Find(".info__small span").First().Text()
	if rawDate == "" {
		return &Article{}, &ArticleError{Index: index, Field: "date", Err: errors.New("missing value")}
	}

	// Витягуємо дату з тексту
	matches := dateRegex.FindStringSubmatch(rawDate)
	if len(matches) == 0 {
		return &Article{}, &ArticleError{Index: index, Field: "date", Err: errors.New("invalid date format")}
	}

	// Перша частина регулярного виразу містить дату
	date := matches[0]

	return &Article{
		Title:   title,
		Link:    link,
		Content: content,
		Image:   imgSrc,
		Date:    date,
	}, nil
}

// FetchLatestArticles об'єднує процес отримання HTML і парсингу статей
func FetchLatestArticles(url string) ([]*Article, error) {
	doc, err := FetchHTML(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching articles: %w", err)
	}

	articles, err := ParseArticle(doc)
	if err != nil {
		return nil, fmt.Errorf("error parsing articles: %w", err)
	}

	return articles, nil
}
