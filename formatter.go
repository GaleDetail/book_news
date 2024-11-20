package main

import "fmt"

func CreatePost(article *Article) string {
	return fmt.Sprintf(
		"📚 *%s*\n"+
			"🗓 *Дата:* %s\n\n"+
			"_%s_\n\n"+
			"🔗 [Читати більше](%s)",
		article.Title, article.Date, article.Content, article.Link,
	)
}
