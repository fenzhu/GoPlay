package page

import (
	"context"
	"errors"
	"log"
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
	select {
	case articleChan <- p:
		//just add article to channel
	default:
		//channel is full now
		go tryBatch(batchSize)
		//block until channel space is available
		articleChan <- p
	}
	return nil
}

func init() {
	go batchWorker(context.Background())
}

var (
	articleChan   = make(chan *Article, batchSize)
	batchInterval = 10 * time.Second
	batchSize     = 100
)

func batchWorker(ctx context.Context) {
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go tryBatch(0)
		}
	}
}

func tryBatch(trigger int) {

	// articles := make([]*Article, 0, len(articleChan))

	// for article := range articleChan {
	// 	articles = append(articles, article)
	// }

	if len(articleChan) > trigger {
		err := db().Transaction(func(tx *gorm.DB) error {
			for article := range articleChan {
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

		// articleChan = articleChan[:0]
	}
}
