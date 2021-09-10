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
	tojCrawler := crawler.TojCrawler()
	uvaCrawler := crawler.UvaCrawler()

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
				switch judgeNmae {
				case "toj":
					url := fmt.Sprintf("https://toj.tfcis.org/oj/be/chal?off=0&proid=%s&acctid=%s", problemId, judges["toj"])
					tojCrawler.Visit(url)
				case "uva":
					ctx := colly.NewContext()
					url := fmt.Sprintf("https://uhunt.onlinejudge.org/api/p/num/%s", problemId)

					ctx.Put("pnum", problemId)
					ctx.Put("uid", judges["uva"])
					uvaCrawler.Request("GET", url, nil, ctx, nil)
				}
			}
		}
	}

	// get all submission
	tojCrawler.Wait()
	uvaCrawler.Wait()
}
