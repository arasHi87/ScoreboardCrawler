package crawler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/arasHi87/ScoreboardCrawler/src/util"
	"github.com/gocolly/colly"
)

type tiojUser struct {
	Id          int    `json:"id"`
	Result      string `json:"result"`
	Score       string `json:"score"`
	ProblemId   int    `json:"problem_id"`
	Submitter   string `json:"submitter"`
	TotalTime   int    `json:"total_time"`
	TotalMemory int    `json:"total_memory"`
	Compiler    string `json:"compiler"`
	CodeLength  int    `json:"code_length"`
	CreatedAt   string `json:"created_at"`
	Url         string `json:"url"`
}

func TiojCrawler() *colly.Collector {
	ctx := context.Background()
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		rdb := util.GetRedis()
		results := make([]tiojUser, 0)
		resp := string(r.Body)

		if err := json.Unmarshal([]byte(resp), &results); err != nil {
			fmt.Println(err.Error())
		}

		pid, uid := r.Ctx.Get("pid"), r.Ctx.Get("uid")
		index := fmt.Sprintf("tioj:%s-%s", pid, uid)

		// insert data into db
		err := rdb.Set(ctx, index, results[0].Result, 0).Err()
		if err != nil {
			panic(err)
		}
	})

	return c
}
