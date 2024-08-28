package page

import (
	"fmt"

	"example.com/wiki/database"
	"gorm.io/gorm"
)

type Page struct {
	Title string
	Body  string
}

func db() *gorm.DB {
	return database.Center.GetDatabase("wiki")
}

func cache() *database.Cache {
	return database.Center.GetCache("wiki")
}

func LoadPage(title string) (*Page, error) {
	page := &Page{Title: title}

	var cache = cache().Data
	body, ok := cache[title]
	if !ok {
		// row := db().QueryRow("SELECT * FROM article WHERE title = ?", title)
		db().Table("article").First(&page, "title = ?", title)
		// if err := row.Scan(&page.Title, &page.Body); err != nil {
		// 	if err == sql.ErrNoRows {
		// 		return &page, fmt.Errorf("pageByTitle %s, no such page", title)
		// 	} else {
		// 		return &page, fmt.Errorf("pageByTitle %s, %v", title, err)
		// 	}
		// }
		fmt.Printf("miss cache %s\n", title)
		if cache[title] != page.Body {
			cache[title] = page.Body
		}
	} else {
		fmt.Printf("hit cache %s\n", title)
		page.Title = title
		page.Body = body
	}

	return page, nil
}

func (p *Page) Save() error {
	_, err := LoadPage(p.Title)

	// var sqlScript string
	// var result sql.Result
	if err == nil {
		// sqlScript = "UPDATE article set body = ? WHERE title = ?"
		// result, err = db().Exec(sqlScript, p.Body, p.Title)
		db().Table("article").Model(p).Where("title = ?", p.Title).Update("body", p.Body)
	} else {
		// sqlScript = "INSERT INTO article (title, body) VALUES (?, ?)"
		// result, err = db().Exec(sqlScript, p.Title, p.Body)
		db().Table("article").Create(p)
	}

	return nil
}
