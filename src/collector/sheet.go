package collector

import (
	"context"
	"fmt"
	"sync"

	"github.com/arasHi87/ScoreboardCrawler/src/util"
)

func SheetCollector(judgeName string, urls []UrlElement, wg *sync.WaitGroup) {
	defer wg.Done()
	submission := util.GetValue(judgeName)
	nameMap := map[string]string{"Kattis": "kattis", "ZeroJudge": "zoj"}
	rdb := util.GetRedis()
	ctx := context.Background()

	for _, url := range urls {
		var result string
		pid := url.Pid
		uid := url.Uid
		key := fmt.Sprintf("%s:%s-%s", nameMap[judgeName], pid, uid)

		if _, ok := submission[uid][pid]; ok {
			result = submission[uid][pid]
		} else {
			result = "NE"
		}

		if err := rdb.Set(ctx, key, result, 0).Err(); err != nil {
			panic(err)
		}
	}
}
