package page

import (
	"errors"

	"example.com/wiki/database"
	"gorm.io/gorm"
)

type Article struct {
	Title string `gorm:"primarykey"`
	Body  string
}

func (p *Article) TableName() string {
	return "article"
}

func db() *gorm.DB {
	return database.Center.GetDatabase("wiki")
}

func cache() *database.Cache {
	return database.Center.GetCache("wiki")
}

func LoadPage(title string) (*Article, error) {
	page := &Article{Title: title}

	var cache = cache().Data
	body, ok := cache[title]
	if !ok {
		res := db().First(&page)
		if res.Error != nil {
			return nil, errors.New("no article found")
		}

		cache[title] = page.Body
	} else {
		page.Title = title
		page.Body = body
	}

	return page, nil
}

func (p *Article) Save() error {
	_, err := LoadPage(p.Title)

	if err == nil {
		return db().Model(p).Update("body", p.Body).Error
	} else {
		return db().Create(p).Error
	}
}
