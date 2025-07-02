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

	// var cache = cache().Data
	// res := cache.Get(context.Background(), title)
	// if res.Err() != nil {
	res := db().First(&page)
	if res.Error != nil {
		return nil, errors.New("no article found")
	}

	// cache.Set(context.Background(), title, page.Body, 0)
	// } else {
	// 	page.Body = res.Val()
	// }

	return page, nil
}

func (p *Article) Save() error {
	// select {
	// case articleChan <- p:
	//just add article to channel
	// default:
	//channel is full now
	// go tryBatch(batchSize)
	//block until channel space is available
	// 	articleChan <- p
	// }

	// update cache
	var cache = cache().Data
	// res := cache.Get(context.Background(), p.Title)
	// if res.Err() != nil {
	cache.Set(context.Background(), p.Title, p.Body, 0)
	// }

	return nil
}

func init() {
	go batchWorker(context.Background())
}

var (
	articleChan   = make(chan *Article, batchSize+1)
	batchInterval = 10 * time.Second
	batchSize     = 1000
	workerSize    = 50
	batchPool     = make(chan struct{}, workerSize)
)

func batchWorker(ctx context.Context) {
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// go tryBatch(0)
		}
	}
}

func tryBatch(trigger int) {
	select {
	case batchPool <- struct{}{}:
		log.Println("batch worker start")
		defer func() {
			log.Println("batch worker end")
			<-batchPool
		}()
		if len(articleChan) > trigger {
			log.Println("batch size ", len(articleChan))
			err := db().Transaction(func(tx *gorm.DB) error {
				for i := 0; i < len(articleChan); i++ {
					article := <-articleChan
					result := db().Clauses(clause.OnConflict{
						Columns:   []clause.Column{{Name: "title"}},
						DoUpdates: clause.AssignmentColumns([]string{"body"}),
					}).Create(article)

					if result.Error != nil {
						return result.Error
					}
				}

				log.Println("transaction success")
				return nil
			})

			if err != nil {
				log.Println("batch write error:", err)
			} else {
				log.Println("batch write success")
			}
		}
	default:
		log.Println("batch work is busy, skip batch")
	}
}
