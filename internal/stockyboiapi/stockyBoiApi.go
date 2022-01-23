package stockyboiapi

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var slackApiKey string
var slackEndpoint string
var slackChannel string

type SlackChatPostMessageText struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

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
	req, _ := http.NewRequest("POST", reqUrl, bytes.NewBuffer(jsonStr))

	req.Header.Add("Authorization", "Bearer "+slackApiKey)
	req.Header.Add("Content-Type", "application/json")

	_, resError := http.DefaultClient.Do(req)
	if resError != nil {
		log.Fatal(resError)
	}
}
