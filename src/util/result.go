package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type ResultSet struct {
	Week map[string][]problemStatus `json:"results"`
}

type problemStatus struct {
	JudgeName string       `json:"judge_name"`
	ProblemId string       `json:"problem_id"`
	Users     []userStatus `json:"users"`
}

type userStatus struct {
	ProblemStatus string `json:"problem_status"`
	UserName      string `json:"user_name"`
}

func IntegrationReseult(homeworkFile []string, users map[string]map[string]string) {
	results := ResultSet{}
	results.Week = make(map[string][]problemStatus)

	// integration result into data struct
	for _, path := range homeworkFile {
		homework := GetHomeworks(path)
		weekCount := strings.Replace(filepath.Base(path), ".json", "", 1)

		for _, problem := range homework.Problems {
			result := problemStatus{
				JudgeName: problem.JudgeName,
				ProblemId: problem.ProblemId,
			}

			for name, judges := range users {
				user := userStatus{UserName: name}
				index := fmt.Sprintf("%s:%s-%s",
					result.JudgeName,
					result.ProblemId,
					judges[result.JudgeName])

				user.ProblemStatus = getStatus(index)
				result.Users = append(result.Users, user)
			}

			results.Week[weekCount] = append(results.Week[weekCount], result)
		}
	}

	// write results into result.json
	file, _ := json.MarshalIndent(results, "", "  ")
	_ = ioutil.WriteFile("data/result.json", file, 0644)
}

func getStatus(index string) string {
	rdb := GetRedis()
	ctx := context.Background()

	status, err := rdb.Get(ctx, index).Result()
	if err != nil {
		panic(err)
	}

	return status
}
