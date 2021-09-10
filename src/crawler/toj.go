package crawler

import (
	"context"
	"fmt"
	"net/url"

	"github.com/arasHi87/ScoreboardCrawler/src/util"
	"github.com/gocolly/colly"
)

var ctx = context.Background()
var tojStatusCodeMap = map[string]string{
	"Accepted":              "AC",
	"Wrong Answer":          "WA",
	"Runtime Error":         "RE",
	"Time Limit Exceeded":   "TLE",
	"Memory Limit Exceeded": "MLE",
	"Compile Error":         "CE",
}

func TojCrawler() *colly.Collector {
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

		// parse url
		u, err := url.Parse(r.URL.String())
		if err != nil {
			panic(err)
		}
		m, _ := url.ParseQuery(u.RawQuery)

		r.Ctx.Put("pid", m["proid"][0])
		r.Ctx.Put("uid", m["acctid"][0])
	})

	return c
}
