package page

import (
	"database/sql"
	"fmt"

	"example.com/wiki/database"
)

type Page struct {
	Title string
	Body  string
}

func db() *sql.DB {
	return database.Center.GetDatabase("wiki")
}

func LoadPage(title string) (*Page, error) {
	var page Page

	page.Title = title
	row := db().QueryRow("SELECT * FROM article WHERE title = ?", title)
	if err := row.Scan(&page.Title, &page.Body); err != nil {
		if err == sql.ErrNoRows {
			return &page, fmt.Errorf("pageByTitle %s, no such page", title)
		} else {
			return &page, fmt.Errorf("pageByTitle %s, %v", title, err)
		}
	}

	return &page, nil
}

func (p *Page) Save() error {
	_, err := LoadPage(p.Title)

	var sqlScript string
	var result sql.Result
	if err == nil {
		sqlScript = "UPDATE article set body = ? WHERE title = ?"
		result, err = db().Exec(sqlScript, p.Body, p.Title)
	} else {
		sqlScript = "INSERT INTO article (title, body) VALUES (?, ?)"
		result, err = db().Exec(sqlScript, p.Title, p.Body)
	}

	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}
