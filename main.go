package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"time"

	"internal/rapidstocks"
	"internal/stockyboiapi"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

// Stores tickers that will be tracked daily
var tickers []string

type AddTickerReq struct {
	Text        string `form:"text" binding:"required"`
	ResponseUrl string `form:"response_url" binding:"required"`
}

// Init the .env file if not running in production
func init() {
	env := os.Getenv("ENVIRONMENT")
	if env != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func main() {
	rapidstocks.Configure(
		os.Getenv("RAPID_STOCKS_URL"),
		os.Getenv("RAPID_STOCKS_TOKEN"),
	)
	stockyboiapi.Configure(
		os.Getenv("SLACK_API_KEY"),
		os.Getenv("SLACK_API_ENDPOINT"),
		os.Getenv("SLACK_API_CHANNEL"),
	)
	StartCron()
	StartGin()
}

func StartCron() {
	asxTimeZone := os.Getenv("TIMEZONE")
	timeZone, _ := time.LoadLocation(asxTimeZone)

	openingJob := cron.NewWithLocation(timeZone)
	openingJob.AddFunc("5 10 * * * *", func() {
		stockSummary(asxTimeZone)
	})

	middayJob := cron.NewWithLocation(timeZone)
	openingJob.AddFunc("0 13 * * * *", func() {
		stockSummary(asxTimeZone)
	})

	closingJob := cron.NewWithLocation(timeZone)
	openingJob.AddFunc("45 15 * * * *", func() {
		stockSummary(asxTimeZone)
	})
	openingJob.Start()
	middayJob.Start()
	closingJob.Start()
}

func StartGin() {
	port := os.Getenv("PORT")
	env := os.Getenv("ENVIRONMENT")
	var address string

	if env == "" {
		log.Fatal("$ENVIRONMENT must be set")
	}

	if env == "production" || env == "development" {
		if env == "production" {
			if port == "" {
				log.Fatal("$PORT must be set")
			}
			address = ":" + port
		} else {
			address = "localhost:8080"
		}
	} else {
		log.Fatal("$ENVIRONMENT not recognised")
	}

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/ping", respondPong)
	router.POST("/showSummary", showSummary)
	router.POST("/showTickers", showTickers)
	router.POST("/addTicker", addTicker)

	router.Run(address)
}

// respondPong is used to run uptime checks and
// definitely not to keep this app running on a heroku
func respondPong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// Gets details on current tickers and posts them
// to slack.
func stockSummary(timezone string) {
	quotes := rapidstocks.GetStocks(tickers)
	formattedQuoteBlocks := stockyboiapi.FormatQuotes(quotes, timezone)
	jsonResp, _ := json.Marshal(formattedQuoteBlocks)
	stockyboiapi.PostToSlack("", jsonResp)
}

// Simple slash command support for showing current tickers
func showSummary(c *gin.Context) {
	stockSummary(os.Getenv("TIMEZONE"))
}

// Automated cron job for posting ticker details
func cronStockSummary(timezone string) {
	time.LoadLocation(timezone)
	day := time.Now().Weekday()
	// Sunday = 0, Saturday = 6
	if day != 0 && day != 6 {
		stockSummary(timezone)
	}
}

// Returns all tickers currently registered for the session.
func showTickers(c *gin.Context) {
	stockyboiapi.SlashCommandShowTickers(tickers)
	c.String(http.StatusOK, "Tickers Sent")
}

// First checks the ticker against the API then stores it.
func addTicker(c *gin.Context) {
	var request AddTickerReq
	err := c.ShouldBind(&request)
	if err != nil {
		log.Fatal(err)
	}
	valid := rapidstocks.ValidateTicker(request.Text)
	fmt.Println(valid)
	if !valid {
		c.String(http.StatusOK, "Ticker not valid")
		return
	}
	tickers = append(tickers, request.Text)
	c.String(http.StatusOK, "Ticker added")
}
