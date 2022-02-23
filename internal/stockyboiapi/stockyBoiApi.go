package stockyboiapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"internal/rapidstocks"
)

var slackApiKey string
var slackEndpoint string
var slackChannel string

type SlackChatPostMessageText struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type SlackRespond struct {
	Text string `json:"text"`
}

type BlocksRespType struct {
	Channel string   `json:"channel"`
	Blocks  []Blocks `json:"blocks"`
}
type Text struct {
	Type  string `json:"type,omitempty"`
	Text  string `json:"text,omitempty"`
	Emoji bool   `json:"emoji,omitempty"`
}
type Fields struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type Blocks struct {
	Type   string   `json:"type"`
	Text   *Text    `json:"text,omitempty"`
	Fields []Fields `json:"fields,omitempty"`
}

// Creates a reqest and posts to a slack URL/Webhook.
func PostToSlack(url string, body []byte) {
	if url == "" {
		url = slackEndpoint + "/chat.postMessage"
	}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))

	req.Header.Add("Authorization", "Bearer "+slackApiKey)
	req.Header.Add("Content-Type", "application/json")

	_, resError := http.DefaultClient.Do(req)
	if resError != nil {
		log.Fatal(resError)
	}
}

// Configures the module with passed variables.
// Must be called before any other methods are called.
func Configure(
	envSlackApiKey string,
	envSlackEndpoint string,
	envSlackChannel string,
) {
	if envSlackApiKey != "" && envSlackEndpoint != "" && envSlackChannel != "" {
		slackApiKey = envSlackApiKey
		slackEndpoint = envSlackEndpoint
		slackChannel = envSlackChannel
	}
}

// Marshalls passed tickers and posts them to slack.
func SlashCommandShowTickers(tickers []string) {
	reqUrl := slackEndpoint + "/chat.postMessage"
	jsonData := &SlackChatPostMessageText{
		Channel: slackChannel,
		Text:    strings.Join(tickers[:], ","),
	}

	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatalf("could not marshal JSON: %s", err)
	}
	PostToSlack(reqUrl, jsonStr)
}

// Posts a message response to slack using the passed URL.
func SlackPostText(message string, responseUrl string) {
	jsonData := &SlackRespond{
		Text: message,
	}
	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatalf("could not marshal JSON: %s", err)
	}
	PostToSlack(responseUrl, jsonStr)
}

func FormatQuotes(quotes []rapidstocks.RespQuote, timezone string) BlocksRespType {
	blocksStruct := BlocksRespType{
		Channel: slackChannel,
	}
	time.LoadLocation(timezone)
	hour := time.Now().Hour()
	// https://app.slack.com/block-kit-builder/T02SDF83ZLN#%7B%22blocks%22:%5B%7B%22type%22:%22header%22,%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Good%20morning!%20:sunrise_over_mountains:%22,%22emoji%22:true%7D%7D,%7B%22type%22:%22divider%22%7D,%7B%22type%22:%22section%22,%22text%22:%7B%22type%22:%22mrkdwn%22,%22text%22:%22*PLB*%20XYZ%22%7D%7D,%7B%22type%22:%22section%22,%22fields%22:%5B%7B%22type%22:%22mrkdwn%22,%22text%22:%22_Opening%20Price_:%20*xyz*%22%7D,%7B%22type%22:%22mrkdwn%22,%22text%22:%22_Previous%20Close_:%20*xyz*%22%7D,%7B%22type%22:%22mrkdwn%22,%22text%22:%22_Range_:%20*xyz*%22%7D,%7B%22type%22:%22mrkdwn%22,%22text%22:%22_Current%20Price_:%20*xyz*%22%7D%5D%7D%5D%7D
	var message string
	if hour > 0 && hour < 11 {
		message = "Good Morning! :sunrise_over_mountains:"
	} else if hour > 11 && hour < 14 {
		message = "Good Afternoon! :desert:"
	} else {
		message = "Good Evening! :city_sunset:"
	}

	headerBlock := Blocks{
		Type: "header",
		Text: &Text{
			Type:  "plain_text",
			Text:  message,
			Emoji: true,
		},
	}
	blocksStruct.Blocks = append(blocksStruct.Blocks, headerBlock)

	for _, ticker := range quotes {
		dividerBlock := Blocks{
			Type: "divider",
		}
		blocksStruct.Blocks = append(blocksStruct.Blocks, dividerBlock)

		tickerTitleBlock := Blocks{
			Type: "section",
			Text: &Text{
				Type: "mrkdwn",
				Text: fmt.Sprintf("%s - *%s*", ticker.Symbol, ticker.RegularMarketPrice.Fmt),
			},
		}
		blocksStruct.Blocks = append(blocksStruct.Blocks, tickerTitleBlock)

		tickerFieldsBlock := Blocks{
			Type: "section",
			Fields: []Fields{
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("_Opening Price_: *%s*", ticker.RegularMarketOpen.Fmt),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("_Previous Close_: *%s*", ticker.RegularMarketPreviousClose.Fmt),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("_Range_: *%s*", ticker.RegularMarketDayRange.Fmt),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("_Current Price_: *%s*", ticker.RegularMarketPrice.Fmt),
				},
			},
		}
		blocksStruct.Blocks = append(blocksStruct.Blocks, tickerFieldsBlock)
	}

	return blocksStruct
}
