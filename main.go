package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const BTCPercent float64 = 70.0
const ETHPercent float64 = 30.0
const ExRateUrl = "https://api.coinbase.com/v2/exchange-rates?currency=USD"

func getExchangeRates() (float64, float64, error) {
	resp, err := http.Get(ExRateUrl)
	if err != nil {
		return 0, 0, fmt.Errorf("error accessing API: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("error reading response body: %v", err)
	}
	result := struct {
		Data struct {
			Currency string `json:"data"`
			Rates    struct {
				BTC string `json:"BTC"`
				ETH string `json:"ETH"`
			} `json:"rates"`
		} `json:"data"`
	}{}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing json response: %v", err)
	}
	btcRate, err := strconv.ParseFloat(result.Data.Rates.BTC, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("btc rate could not be cast to float64")
	}
	ethRate, err := strconv.ParseFloat(result.Data.Rates.ETH, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("eth rate could not be cast to float64")
	}
	return btcRate, ethRate, nil
}

func allocateInvestment(investAmount float64) (float64, float64) {
	btcUsdAmt := investAmount * BTCPercent / 100
	ethUsdAmt := investAmount * ETHPercent / 100
	return btcUsdAmt, ethUsdAmt
}

func usdToCrypto(amountUSD float64, exRate float64) float64 {
	amountCrypto := amountUSD * exRate
	return amountCrypto
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Please pass the amount you want to invest on the command line")
		return
	}
	amount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		fmt.Println("Please pass only numbers")
		return
	}
	btcRate, ethRate, err := getExchangeRates()
	if err != nil {
		fmt.Println(err)
		return
	}
	btcUsdAmount, ethUsdAmount := allocateInvestment(amount)
	btcAmount := usdToCrypto(btcUsdAmount, btcRate)
	ethAmount := usdToCrypto(ethUsdAmount, ethRate)
	investments := struct {
		BTC string
		ETH string
	}{strconv.FormatFloat(btcAmount, 'f', 8, 64),
		strconv.FormatFloat(ethAmount, 'f', 8, 64)}
	jsonOutput, err := json.Marshal(investments)
	if err != nil {
		fmt.Printf("could not marshal json output: %v\n", err)
		return
	}
	fmt.Println(string(jsonOutput))
}
