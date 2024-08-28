package page

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"example.com/wiki/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	articleMutex.Lock()
	articles = append(articles, p)
	articleMutex.Unlock()

	cache().Data[p.Title] = p.Body
	tryBatch(batchSize)
	return nil
}

func init() {
	go batchWorker()
}

var (
	articles      = make([]*Article, 0, batchSize)
	batchInterval = 10 * time.Second
	batchSize     = 100
	articleMutex  sync.Mutex
)

func batchWorker() {
	for {
		time.Sleep(batchInterval)
		tryBatch(0)
	}
}

func tryBatch(trigger int) {
	if len(articles) > trigger {
		fmt.Printf("start batch, trigger %d\n", trigger)
		err := db().Transaction(func(tx *gorm.DB) error {
			for _, article := range articles {
				result := db().Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "title"}},
					DoUpdates: clause.AssignmentColumns([]string{"body"}),
				}).Create(article)

				if result.Error != nil {
					return result.Error
				}
			}

			return nil
		})

		if err != nil {
			log.Println("batch write error:", err)
		}

		articles = articles[:0]
	}
}
