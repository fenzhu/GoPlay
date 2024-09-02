package service

import (
	"context"
	"strings"
	"time"

	"example.com/vote/database"
	"github.com/redis/go-redis/v9"
)

var (
	ONE_WEEK_IN_SECONDS int64   = 7 * 86400
	VOTE_SCORE          float64 = 432
	ARTICLE_PREFIX      string  = "article:"
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
func article_vote(user string, article string) {
	client := client()

	cutoff := time.Now().Unix() - ONE_WEEK_IN_SECONDS
	publishTime := client.ZScore(context.Background(), "time", article)
	//a week passed, can't vote
	if cutoff > int64(publishTime.Val()) {
		return
	}

	articleId := getArticleID(article)
	//更新投票信息
	suc := client.SAdd(context.Background(), "voted:"+articleId, user).Val()

	if suc == 1 {
		//增加分数
		client.ZIncrBy(context.Background(), "score", VOTE_SCORE, article)
		//更新文章信息
		client.HIncrBy(context.Background(), article, "votes", 1)
	}
}

func article_post(user string, title string, link string) {
	client := client()
	articleId := string(client.Incr(context.Background(), "article:").Val())

	//更新文章信息
	now := time.Now().Unix()
	article := ARTICLE_PREFIX + articleId
	client.HMSet(context.Background(), article, map[string]interface{}{
		"title":  title,
		"link":   link,
		"poster": user,
		"time":   now,
		"votes":  1,
	})

	//更新投票信息
	voted := "voted:" + articleId
	client.SAdd(context.Background(), voted, user)

	//更新分数
	client.ZAdd(context.Background(), "score", redis.Z{Member: article, Score: VOTE_SCORE})
	//更新发布时间
	client.ZAdd(context.Background(), "time", redis.Z{Member: article, Score: float64(now)})
}
