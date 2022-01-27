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

type SlackRespond struct {
	Text string `json:"text"`
}

//Internal: Creates a reqest and posts to a slack URL/Webhook.
func postToSlack(url string, body []byte) {
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
	postToSlack(reqUrl, jsonStr)
}

// Posts a message response to slack using the passed URL.
func SlackRespondToSlashCommand(message string, responseUrl string) {
	jsonData := &SlackRespond{
		Text: message,
	}
	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatalf("could not marshal JSON: %s", err)
	}
	postToSlack(responseUrl, jsonStr)
}
