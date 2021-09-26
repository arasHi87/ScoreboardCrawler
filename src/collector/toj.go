package collector

import (
	"context"
	"fmt"
	"sync"

	"github.com/arasHi87/ScoreboardCrawler/src/util"
	"github.com/gocolly/colly"
)

var tojStatusCodeMap = map[string]string{
	"Accepted":              "AC",
	"Wrong Answer":          "WA",
	"Runtime Error":         "RE",
	"Time Limit Exceeded":   "TLE",
	"Memory Limit Exceeded": "MLE",
	"Compile Error":         "CE",
}

func TojCollector(urls []UrlElement, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx := context.Background()
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)

	c.OnHTML("table#challist>tbody", func(e *colly.HTMLElement) {
		rdb := util.GetRedis()
		status := e.ChildText("td.state-1")

		pid, uid := e.Request.Ctx.Get("pid"), e.Request.Ctx.Get("uid")
		index := fmt.Sprintf("toj:%s-%s", pid, uid)

		// convert status code
		var problemStatus string
		if status == "" {
			problemStatus = "NE"
		} else if val, ok := tojStatusCodeMap[e.Text]; ok {
			problemStatus = val
		} else {
			problemStatus = "OE"
		}

		// insert data into db
		err := rdb.Set(ctx, index, problemStatus, 0).Err()
		if err != nil {
			panic(err)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	for _, url := range urls {
		ctx := colly.NewContext()

		ctx.Put("pid", url.Pid)
		ctx.Put("uid", url.Uid)
		c.Request("GET", url.Url, nil, ctx, nil)
	}
	c.Wait()
}
