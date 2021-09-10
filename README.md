# SCIST ScordBoard crawler

A crawler made by colly to get all user submissions on every judge.

## Development

### Prerequisitive

| Name | Version |
| ---- | ------- |
| go   | 1.17    |

### Environment setup

0. Download required packages

```
go mod download
```

1. Start run crawler

```
cd src/
go run main.go
```

### Formatting

This project use `gofmt` for formatting

```
make format
```

## Data

This Project support following judge, and will read `homework`, `user` file to build request.

- [TOJ](https://toj.tfcis.org)
- [UVA](https://onlinejudge.org)

### Structure

0. Homework file in `src/data/problems/`, and make sure that the file name format complies with `week*.json`.

```json
{
  "problems": [
    {
      "judge_name": "uva",
      "problem_id": "1595"
    },
    {
      "judge_name": "toj",
      "problem_id": "87"
    }
  ]
}
```

1. User file in `src/data/user.json`.

```json
{
  "arashi87": {
    "toj": "1931",
    "uva": "1002251"
  }
}
```
