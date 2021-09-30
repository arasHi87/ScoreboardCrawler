package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tok, err := tokenFromEnv()
	if err != nil {
		panic("Token expired, please re-generate")
	}
	return config.Client(context.Background(), tok)
}

// Retrieves a token from a env.
func tokenFromEnv() (*oauth2.Token, error) {
	env := GetEnv()
	tok := &oauth2.Token{}
	err := json.NewDecoder(strings.NewReader(env["TOKEN"])).Decode(tok)
	return tok, err
}

// Retrieves value from spreadsheet
func GetValue(sheetName string) map[string]map[string]string {
	fmt.Println("Visiting", sheetName)
	submissions := make(map[string]map[string]string)
	ctx := context.Background()
	env := GetEnv()
	b := []byte(env["CREDENTIALS"])

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId := env["SHEETID"]
	readRange := sheetName
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for i := 1; i < len(resp.Values); i++ {
			title := resp.Values[0] // use to get pid
			value := resp.Values[i]
			userName := fmt.Sprintf("%v", value[1])
			for j := 3; j < len(value); j += 2 {
				pid := fmt.Sprintf("%v", title[j])
				result := fmt.Sprintf("%v", value[j])
				if _, ok := submissions[userName]; !ok {
					submissions[userName] = make(map[string]string)
				}
				submissions[userName][pid] = result
			}
		}
	}

	return submissions
}
