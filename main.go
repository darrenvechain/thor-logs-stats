package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Event struct {
	LogStreamName string `json:"logStreamName"`
	Timestamp     int64  `json:"timestamp"`
	Message       string `json:"message"`
	IngestionTime int64  `json:"ingestionTime"`
	EventId       string `json:"eventId"`
}

type Logs struct {
	Events []Event `json:"events"`
}

type EventOptions struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type RequestBody struct {
	Options EventOptions `json:"options"`
}

func main() {
	file, err := os.Open("logs_output.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var logs Logs

	// Decode the JSON data
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&logs); err != nil {
		panic(err)
	}

	optionsArray := make([]EventOptions, 0)

	for _, event := range logs.Events {
		if strings.Contains(event.Message, "URI=/logs/event") || strings.Contains(event.Message, "URI=/logs/transfer") {
			// Regular expression to extract the JSON Body
			re := regexp.MustCompile(`Body="(\{.*?\})"`)
			matches := re.FindStringSubmatch(event.Message)
			if len(matches) > 1 {
				reqBody := RequestBody{
					Options: EventOptions{
						Offset: -1,
						Limit:  -1,
					},
				}
				// Parse the JSON Body into EventOptions
				cleanedJSON := strings.Replace(matches[1], `\"`, `"`, -1)
				cleanedJSON = strings.Replace(cleanedJSON, `\\`, `\`, -1)
				err := json.Unmarshal([]byte(cleanedJSON), &reqBody)
				if err != nil {
					panic(err)
				}

				optionsArray = append(optionsArray, reqBody.Options)
			}
		}
	}

	var noLimitCount int
	var greaterThan1000Count int
	for _, options := range optionsArray {
		if options.Limit == -1 {
			noLimitCount++
		} else if options.Limit > 1000 {
			greaterThan1000Count++
		}
	}

	fmt.Println("Total number of logs queries: ", len(optionsArray))
	fmt.Println("Number of queries with no limit: ", noLimitCount)
	fmt.Println("Number of queries with limit greater than 1000: ", greaterThan1000Count)
}
