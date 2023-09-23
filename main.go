package main

import (
	"fmt"
	"log"

	"crypto/hmac"
	"crypto/sha256"
	"io"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

type Binance struct {
	API_KEY    string
	API_SECRET string
	BASE_URL   string
}

func (b *Binance) router(endpoint string, params string, method string) string {
	url := b.BASE_URL + endpoint + params + "&signature=" + b.signature(params)

	fmt.Println(url)
	client := &http.Client{}

	getReq, _ := http.NewRequest("GET", url, nil)
	getReq.Header.Add("X-MBX-APIKEY", b.API_KEY)

	resp, _ := client.Do(getReq)
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return string(res)
}

func (b *Binance) signature(request string) string { // hmac sha256 signature
	sign := hmac.New(sha256.New, []byte(b.API_SECRET))
	sign.Write([]byte(request))
	return fmt.Sprintf("%x", sign.Sum(nil))
}

func (b *Binance) getTimeStamp() int64 {
	return time.Now().Unix() * 1000
}

func (b *Binance) getFuturesBalance() float64 {
	rawJson := b.router("/fapi/v2/balance?", fmt.Sprintf("timestamp=%d", b.getTimeStamp()), "get")
	balance := gjson.Get(rawJson, "5.balance")
	crossUnPnl := gjson.Get(rawJson, "5.crossUnPnl")

	return balance.Float() + crossUnPnl.Float()
}

func main() {
	myAcc := Binance{}

	fmt.Println(myAcc.getFuturesBalance())
}
