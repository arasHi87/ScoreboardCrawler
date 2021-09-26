package collector

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

func TiojCollector(urls []UrlElement) {

	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		var status string
		rdb := util.GetRedis()
		ctx := context.Background()
		results := make([]tiojUser, 0)
		resp := string(r.Body)

		if err := json.Unmarshal([]byte(resp), &results); err != nil {
			fmt.Println(err.Error())
		}

		pid, uid := r.Ctx.Get("pid"), r.Ctx.Get("uid")
		index := fmt.Sprintf("tioj:%s-%s", pid, uid)

		// map status
		if results[0].Result == "IE" {
			status = "OE"
		} else {
			status = results[0].Result
		}

		// insert data into db
		err := rdb.Set(ctx, index, status, 0).Err()
		if err != nil {
			panic(err)
		}
	})

	for _, url := range urls {
		ctx := colly.NewContext()

		ctx.Put("pid", url.Pid)
		ctx.Put("uid", url.Uid)
		c.Request("GET", url.Url, nil, ctx, nil)
	}
	c.Wait()
}
