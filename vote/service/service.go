package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"example.com/vote/database"
	"github.com/redis/go-redis/v9"
)

var (
	ONE_WEEK_IN_SECONDS int64   = 7 * 86400
	VOTE_SCORE          float64 = 432
	ARTICLE_PREFIX      string  = "article:"
	ARTICLES_PER_PAGE   int64   = 25
)

func client() *redis.Client {
	return database.GetRedis()
}

func getArticleID(article string) string {
	return strings.TrimPrefix(article, ARTICLE_PREFIX)
}

/*
time 文章发布时间
score 文章分数

voted:(文章id) 文章投票用户
article:(文章id) 文章信息
*/
func ArtivleVote(user string, article string) error {
	client := client()

	cutoff := time.Now().Unix() - ONE_WEEK_IN_SECONDS
	publishTime := client.ZScore(context.Background(), "time", article)
	//a week passed, can't vote
	if cutoff > int64(publishTime.Val()) {
		return errors.New("a week passed, can't vote")
	}

	articleId := getArticleID(article)
	//更新投票信息
	suc := client.SAdd(context.Background(), "voted:"+articleId, user).Val()

	if suc == 1 {
		//增加分数
		client.ZIncrBy(context.Background(), "score", VOTE_SCORE, article)
		//更新文章信息
		client.HIncrBy(context.Background(), article, "votes", 1)
		return errors.New("user already voted")
	}
	return nil
}

func ArticlePost(user string, title string, link string) string {
	client := client()
	res := client.Incr(context.Background(), "article:")
	if res.Err() != nil {
		return res.Err().Error()
	}
	id := res.Val()
	fmt.Printf("next id: %d\n", id)
	articleId := fmt.Sprintf("%d", id)
	// strings.i
	//更新文章信息
	now := time.Now().Unix()
	article := ARTICLE_PREFIX + articleId
	client.HMSet(context.Background(), article, map[string]interface{}{
		"title":    title,
		"link":     link,
		"poster":   user,
		"posttime": now,
		"votes":    1,
	})

	//更新投票信息
	voted := "voted:" + articleId
	client.SAdd(context.Background(), voted, user)

	//更新分数
	client.ZAdd(context.Background(), "score", redis.Z{Member: article, Score: VOTE_SCORE})
	//更新发布时间
	client.ZAdd(context.Background(), "time", redis.Z{Member: article, Score: float64(now)})

	return articleId
}

func GetArticles(page int64) []map[string]string {
	start := (page - 1) * ARTICLES_PER_PAGE
	end := start + ARTICLES_PER_PAGE - 1

	client := client()
	articleIds := client.ZRevRange(context.Background(), "score", start, end).Val()

	articles := make([]map[string]string, 0, len(articleIds))
	for _, ararticleId := range articleIds {
		article := client.HGetAll(context.Background(), ARTICLE_PREFIX+ararticleId).Val()

		articles = append(articles, article)
	}
	return articles
}

func GetArticle(id string) map[string]string {
	client := client()
	article := client.HGetAll(context.Background(), ARTICLE_PREFIX+id).Val()
	return article
}
