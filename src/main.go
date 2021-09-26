package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/arasHi87/ScoreboardCrawler/src/collector"
	"github.com/arasHi87/ScoreboardCrawler/src/util"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	var users map[string]map[string]string
	urls := map[string][]collector.UrlElement{
		"toj":  make([]collector.UrlElement, 0),
		"uva":  make([]collector.UrlElement, 0),
		"tioj": make([]collector.UrlElement, 0),
	}

	// load users
	file, _ := ioutil.ReadFile("data/user.json")
	_ = json.Unmarshal([]byte(file), &users)

	// load homeworks
	homeworkFile, err := filepath.Glob("data/problems/week*.json")
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range homeworkFile {
		homework := util.GetHomeworks(path)
		for _, problem := range homework.Problems {
			judgeNmae := problem.JudgeName
			problemId := problem.ProblemId

			for _, judges := range users {
				var url string
				userId := judges[judgeNmae]
				element := collector.UrlElement{
					Pid: problemId,
					Uid: userId,
				}

				// generate url
				switch judgeNmae {
				case "toj":
					url = fmt.Sprintf("https://toj.tfcis.org/oj/be/chal?off=0&proid=%s&acctid=%s", problemId, userId)
				case "tioj":
					url = fmt.Sprintf("https://tioj.ck.tp.edu.tw/submissions.json?filter_username=%s&filter_problem=%s", userId, problemId)
				}

				element.Url = url
				urls[judgeNmae] = append(urls[judgeNmae], element)
			}
		}
	}

	// run collector
	for judgeNmae, urls := range urls {
		switch judgeNmae {
		case "toj":
			collector.TojCollector(urls)
		case "uva":
			collector.UvaCollector(urls)
		case "tioj":
			collector.TiojCollector(urls)
		}
	}

	// integration all result into result.json
	util.IntegrationReseult(homeworkFile, users)
}
