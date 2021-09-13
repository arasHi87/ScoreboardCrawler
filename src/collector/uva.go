package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

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

func UvaCollector() *colly.Collector {
	ctx := context.Background()
	numCollector := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)
	statusCollector := numCollector.Clone()

	numCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	numCollector.OnResponse(func(r *colly.Response) {
		uid := r.Ctx.Get("uid")
		problem := make(map[string]int)
		resp := string(r.Body)

		json.Unmarshal([]byte(resp), &problem)
		r.Ctx.Put("pid", problem["pid"])

		url := fmt.Sprintf("https://uhunt.onlinejudge.org/api/subs-pids/%s/%d/0", uid, problem["pid"])
		statusCollector.Request("GET", url, nil, r.Ctx, nil)
		statusCollector.Wait()
	})

	statusCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	statusCollector.OnResponse(func(r *colly.Response) {
		rdb := util.GetRedis()
		result := make(map[string]uvaUser)
		resp := string(r.Body)

		if err := json.Unmarshal([]byte(resp), &result); err != nil {
			fmt.Println(err.Error())
		}

		// conver status code
		for userId, val := range result {
			problemNum := r.Ctx.Get("pnum")
			index := fmt.Sprintf("uva:%s-%s", problemNum, userId)

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

	return numCollector
}
