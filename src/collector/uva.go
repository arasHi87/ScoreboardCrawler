package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/arasHi87/ScoreboardCrawler/src/util"
	"github.com/gocolly/colly"
)

var uvaStatusCodeMap = map[string]string{
	"90": "AC",
	"70": "WA",
	"40": "RE",
	"50": "TLE",
	"60": "MLE",
	"30": "CE",
}

type uvaUser struct {
	Name        string  `json:"name"`
	UName       string  `json:"uname"`
	Submissions [][]int `json:"subs"`
}

func UvaCollector(urls []UrlElement) {
	// !important thing
	// In order to maintain the unity of the data structure, the pids here are pnum
	uids := ""
	pids := make(map[string]bool)
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		result := make(map[string]uvaUser)
		resp := string(r.Body)

		if err := json.Unmarshal([]byte(resp), &result); err != nil {
			fmt.Println(err.Error())
		}

		for uid, val := range result {
			rdb := util.GetRedis()
			pnum := r.Ctx.Get("pnum")
			ctx := context.Background()
			index := fmt.Sprintf("uva:%s-%s", pnum, uid)

			if len(val.Submissions) > 0 {
				verdictId := strconv.Itoa(val.Submissions[0][2])

				// convert verdict id to problemStatus
				var problemStatus string
				if code, ok := uvaStatusCodeMap[verdictId]; ok {
					problemStatus = code
				} else {
					problemStatus = "OE"
				}

				// insert result into db
				if err := rdb.Set(ctx, index, problemStatus, 0).Err(); err != nil {
					panic(err)
				}
			} else {
				if err := rdb.Set(ctx, index, "NE", 0).Err(); err != nil {
					panic(err)
				}
			}
		}
	})

	// classification pid and uid
	for _, url := range urls {
		uid := url.Uid
		pids[url.Pid] = true

		// insert pid and uid into string
		if i := strings.Index(uids, uid); i == -1 {
			uids = uids + uid + ","
		}
	}

	// build url string and pass to collector
	for pid := range pids {
		ctx := colly.NewContext()
		url := fmt.Sprintf("https://uhunt.onlinejudge.org/api/subs-nums/%s/%s/0", uids, pid)

		ctx.Put("pnum", pid) // watch out pid here is pnum
		c.Request("GET", url, nil, ctx, nil)
	}

	c.Wait()
}
