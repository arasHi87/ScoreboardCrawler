package util

import (
	"encoding/json"
	"io/ioutil"
)

type Homework struct {
	Problems []Problem `json:"problems"`
}

type Problem struct {
	JudgeName string `json:"judge_name"`
	ProblemId string `json:"problem_id"`
}

func GetHomeworks(filename string) Homework {
	file, _ := ioutil.ReadFile(filename)
	data := Homework{}

	_ = json.Unmarshal([]byte(file), &data)
	return data
}
