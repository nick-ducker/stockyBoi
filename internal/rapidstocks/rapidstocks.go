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

func Configure(envUrl string, envToken string) {
	url = envUrl
	token = envToken
	if url == "" || token == "" {
		log.Fatal("Url and Token not configured")
	}
}

//Get one stock
func GetStock(ticker string) []RespQuote {
	reqUrl := "https://" + url + reqAppend + "/quote?symbols=" + ticker
	res := getRequest(reqUrl)

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var body GetStockResponse
	json.Unmarshal(bodyBytes, &body)

	fmt.Println(res)
	fmt.Println(body)

	return body.Data.QuoteResponse.Result

}

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

func Task() {
	fmt.Println("I am running task.")
}
