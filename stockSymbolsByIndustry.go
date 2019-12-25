package main

import (
	"fmt"
	"regexp"
	"net/http"
	"io/ioutil"
	"strings"
	"encoding/json"
	"time"
	"sync"
)

type StockSymbol struct{
	Code string `json:"symbol_code"`
	Description string `json:"company_name"`
}

 const (
 	Url string = "http://bigcharts.marketwatch.com/industry/bigcharts-com/stocklist.asp?Symb=%s&startingIndex=%d"
 	Pattern string = "<td class=\"symb-col\">[A-Za-z0-9.]*</td>\\s*<td class=\"name-col\"><div>.*</div>"
 	Agriculture string = "WSJMXUSAGRI"
 )

func getStockSymbols(industry string, page int, wg *sync.WaitGroup) {
	start := time.Now()
 	Symbols := []StockSymbol{}
 	urlToRead := fmt.Sprintf(Url, industry, page)
 	resp, err := http.Get(urlToRead)
 	if err != nil {
 		fmt.Printf("Could  not read url: %s with error: %v\n", urlToRead, err)
 	}
 	b, err := ioutil.ReadAll(resp.Body)
 	if err != nil {
 		fmt.Printf("Could not read as bytes: %v\n", err)
 	}
 	re := regexp.MustCompile(Pattern)
	for _, res := range re.FindAll(b, -1) {
 		s := fmt.Sprintf("%s", res)
 		i := strings.Index(s, ">") + 1
 		j := strings.Index(s, "</td>")
 		x := strings.Index(s, "<div>") + 5
 		y := strings.Index(s, "</div>")
 		Symbols = append(Symbols,StockSymbol{Code: s[i:j], Description: s[x:y]})
	}

 	resp.Body.Close()
 	data, err := json.Marshal(Symbols)
 	if err != nil {
 		fmt.Println("Could not marshal struct to JSON")
 	}
 	elapsed := time.Since(start)
 	fmt.Printf("Number of symbols %d time of: %s page: %d \n", len(data), elapsed, page)
 	wg.Done()
}

 func main() {
 	mainStart := time.Now()
	
	defer func() {
		mainElapsed := time.Since(mainStart)
		fmt.Printf("\nTotal time: %s\n", mainElapsed)
	}()

 	pages := []int{0, 50, 100, 150, 200}
 	var wg sync.WaitGroup
 	wg.Add(len(pages))

 	for _, page := range pages {
 		go getStockSymbols(Agriculture, page, &wg)
 	}

 	wg.Wait()
 	// start := time.Now()
 	// Symbols := []StockSymbol{}
 	// page := 0
 	// urlToRead := fmt.Sprintf(Url, Agriculture, page)
 	// resp, err := http.Get(urlToRead)
 	// if err != nil {
 	// 	fmt.Printf("Could  not read url: %s with error: %v\n", urlToRead, err)
 	// }
 	//  b, err := ioutil.ReadAll(resp.Body)
 	// if err != nil {
 	// 	fmt.Printf("Could not read as bytes: %v\n", err)
 	// }
 	// re := regexp.MustCompile(Pattern)
 	// results := re.FindAll(b, -1)

 	// for len(results) > 0 {
 	// 	for _, res := range results {
	 // 		s := fmt.Sprintf("%s", res)
	 // 		i := strings.Index(s, ">") + 1
	 // 		j := strings.Index(s, "</td>")
	 // 		x := strings.Index(s, "<div>") + 5
	 // 		y := strings.Index(s, "</div>")
	 // 		Symbols = append(Symbols,StockSymbol{Code: s[i:j], Description: s[x:y]})
 	// 	}

 	// 	page += 50
 	// 	urlToRead = fmt.Sprintf(Url, Agriculture, page)
 	// 	// fmt.Println(urlToRead)
 	// 	resp, err = http.Get(urlToRead)
 	// 	if err != nil {
 	// 		fmt.Printf("Could  not read url: %s with error: %v\n", urlToRead, err)
 	// 	}

 	// 	b, err := ioutil.ReadAll(resp.Body)
 	// 	if err != nil {
 	// 		fmt.Printf("Could not read as bytes: %v\n", err)
 	// 	}
 	// 	results = re.FindAll(b, -1)
 	// 	// fmt.Println(resp.Status)
 	// }
 	// resp.Body.Close()
 	// data, err := json.Marshal(Symbols)
 	// if err != nil {
 	// 	fmt.Println("Could not marshal struct to JSON")
 	// }
 	// elapsed := time.Since(start)
 	// fmt.Printf("Number of symbols %d too total of: %s \n", len(data), elapsed)

 }