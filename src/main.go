package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/arasHi87/ScoreboardCrawler/src/crawler"
	"github.com/arasHi87/ScoreboardCrawler/src/util"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	var users map[string]map[string]string
	crawlers := map[string]*colly.Collector{
		"toj":  crawler.TojCrawler(),
		"uva":  crawler.UvaCrawler(),
		"tioj": crawler.TiojCrawler(),
	}

	// load users
	file, _ := ioutil.ReadFile("data/user.json")
	_ = json.Unmarshal([]byte(file), &users)

	// load homeworks
	homeworkFile, err := filepath.Glob("data/problems/week*.json")
	if err != nil {
		log.Fatal(err)
	}

	// build url set
	for _, path := range homeworkFile {
		homework := util.GetHomeworks(path)
		for _, problem := range homework.Problems {
			judgeNmae := problem.JudgeName
			problemId := problem.ProblemId

			for _, judges := range users {
				ctx := colly.NewContext()
				userId := judges[judgeNmae]

				ctx.Put("pid", problemId)
				ctx.Put("uid", userId)

				switch judgeNmae {
				case "toj":
					url := fmt.Sprintf("https://toj.tfcis.org/oj/be/chal?off=0&proid=%s&acctid=%s", problemId, userId)
					crawlers[judgeNmae].Request("GET", url, nil, ctx, nil)
				case "uva":
					url := fmt.Sprintf("https://uhunt.onlinejudge.org/api/p/num/%s", problemId)
					ctx.Put("pnum", problemId)
					crawlers[judgeNmae].Request("GET", url, nil, ctx, nil)
				case "tioj":
					url := fmt.Sprintf("https://tioj.ck.tp.edu.tw/submissions.json?filter_username=%s&filter_problem=%s", userId, problemId)
					crawlers[judgeNmae].Request("GET", url, nil, ctx, nil)
				}
			}
		}
	}

	// get all submission
	for _, crawler := range crawlers {
		crawler.Wait()
	}

	// integration all result into result.json
	util.IntegrationReseult(homeworkFile, users)
}
