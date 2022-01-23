package rapidstocks

// uses this boi
// https://rapidapi.com/datascraper/api/live-stock-market/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Status int
	Data   struct {
		QuoteResponse struct {
			Result []RespQuote
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
	req, _ := http.NewRequest("GET", reqUrl, nil)

	req.Header.Add("x-rapidapi-host", url)
	req.Header.Add("x-rapidapi-key", token)

	res, resError := http.DefaultClient.Do(req)
	if resError != nil {
		log.Fatal(resError)
	}

	defer res.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(res.Body)

	var body Response
	json.Unmarshal(bodyBytes, &body)

	fmt.Println(res)
	fmt.Println(body)

	return body.Data.QuoteResponse.Result

}

//Get list of stock

func Task() {
	fmt.Println("I am running task.")
}
