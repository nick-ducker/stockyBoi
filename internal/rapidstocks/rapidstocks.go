package rapidstocks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type GetStockResponse struct {
	Status int
	Data   struct {
		QuoteResponse struct {
			Result []RespQuote
		}
	}
}

type ValidateTickerResponse struct {
	Status int
	Data   struct {
		SymbolsValidation struct {
			Result []map[string]bool
		}
	}
}

type RespQuote struct {
	Symbol            string
	RegularMarketOpen struct {
		Raw float64
		Fmt string
	}
	RegularMarketDayRange struct {
		Raw string
		Fmt string
	}
	RegularMarketPrice struct {
		Raw float64
		Fmt string
	}
	RegularMarketPreviousClose struct {
		Raw float64
		Fmt string
	}
}

var url string
var token string

const reqAppend string = "/yahoo-finance/v1"

// Performs a get request to the rapidstocks api.
// Requires the module to be configured with the Configure method.
func getRequest(reqUrl string) *http.Response {
	req, _ := http.NewRequest("GET", reqUrl, nil)

	req.Header.Add("x-rapidapi-host", url)
	req.Header.Add("x-rapidapi-key", token)

	res, resError := http.DefaultClient.Do(req)
	if resError != nil {
		log.Fatal(resError.Error())
	}

	return res
}

// Configures the module with passed variables.
// Must be called before any other methods are called.
func Configure(envUrl string, envToken string) {
	url = envUrl
	token = envToken
	if url == "" || token == "" {
		log.Fatal("Url and Token not configured")
	}
}

/**
Queries the RapidAPI stocks endpoint and returns.
=> Regular market open
=> Regular market day range
=> Regular market price
=> Regular market previous close
**/
func GetStocks(tickers []string) []RespQuote {
	formattedTickers := strings.Join(tickers, "%2")
	reqUrl := "https://" + url + reqAppend + "/quote?symbols=" + formattedTickers
	res := getRequest(reqUrl)

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var body GetStockResponse
	json.Unmarshal(bodyBytes, &body)

	return body.Data.QuoteResponse.Result
}

// Takes one ticker and checks it against RapidAPI.
func ValidateTicker(ticker string) bool {
	strippedString := strings.ReplaceAll(ticker, " ", "")
	reqUrl := "https://" + url + reqAppend + "/symbols-validation?symbols=" + strippedString
	res := getRequest(reqUrl)
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var body ValidateTickerResponse
	json.Unmarshal(bodyBytes, &body)
	firstTicker := body.Data.SymbolsValidation.Result[0]
	fmt.Println(body.Data.SymbolsValidation.Result)
	if firstTicker[ticker] == true {
		return true
	}
	return false
}

// Junk for now.
func Task() {
	fmt.Println("I am running task.")
}
