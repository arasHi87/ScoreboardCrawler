package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/arasHi87/ScoreboardCrawler/src/util"
	"github.com/gocolly/colly"
)

type atcoderSubmission struct {
	UserId      string `json:"user_id"`
	Result      string `json:"result"`
	ProblemId   string `json:"problem_id"`
	EpochSecond int    `json:"epoch_second"`
}

var atcoderStatusCodeMap = map[string]string{
	"AC":  "AC",
	"WA":  "WA",
	"RE":  "RE",
	"TLE": "TLE",
	"MLE": "MLE",
	"CE":  "CE",
}

func AtcoderCollector(urls []UrlElement, wg *sync.WaitGroup) {
	defer wg.Done()
	submissions := make(map[string]string)
	uids, pids := make(map[string]bool), make(map[string]bool)
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		result := make([]atcoderSubmission, 0)
		resp := string(r.Body)

		if err := json.Unmarshal([]byte(resp), &result); err != nil {
			fmt.Println(err.Error())
		}

		// sort result, ac first, time second
		sort.SliceStable(result, func(i, j int) bool {
			if result[i].ProblemId == result[j].ProblemId {
				if result[i].Result == "AC" {
					return true
				}
				if result[j].Result == "AC" {
					return false
				}
				return result[i].EpochSecond < result[j].EpochSecond
			}
			return result[i].ProblemId > result[j].ProblemId
		})

		// map result to target
		for _, submission := range result {
			uid := submission.UserId
			pid := submission.ProblemId
			result := submission.Result

			// insert if value none exist
			if _, ok := submissions[uid+pid]; !ok {

				// map result
				if _, ok := atcoderStatusCodeMap[result]; ok {
					submissions[uid+pid] = result
				} else {
					submissions[uid+pid] = "OE"
				}
			}
		}
	})

	// build url and record target
	for _, url := range urls {
		uid := url.Uid
		pids[url.Pid] = true

		if _, ok := uids[uid]; !ok {
			uids[uid] = true
			url := fmt.Sprintf("https://kenkoooo.com/atcoder/atcoder-api/v3/user/submissions?user=%s&from_second=0", uid)

			c.Visit(url)
		}
	}
	c.Wait()

	// integration result
	rdb := util.GetRedis()
	for uid := range uids {
		for pid := range pids {
			result := ""
			index := uid + pid
			ctx := context.Background()
			key := fmt.Sprintf("atcoder:%s-%s", pid, uid)

			if _, ok := submissions[index]; ok {
				result = submissions[index]
			} else {
				result = "NE"
			}

			if err := rdb.Set(ctx, key, result, 0).Err(); err != nil {
				panic(err)
			}
		}
	}
}
