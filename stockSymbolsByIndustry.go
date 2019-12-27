package main

import (
	"fmt"
	"regexp"
	"net/http"
	"strings"
	"strconv"
	"encoding/json"
	"time"
	"sync"
	"math"
	
	"os"
	"io/ioutil"
	// "log"
)

type StockSymbol struct{
	Code string `json:"symbol_code"`
	Description string `json:"company_name"`
}

 const (
 	Url string = "http://bigcharts.marketwatch.com/industry/bigcharts-com/stocklist.asp?Symb=%s&startingIndex=%d"
 	StockPattern string = "<td class=\"symb-col\">[A-Za-z0-9.]*</td>\\s*<td class=\"name-col\"><div>.*</div>"
 	PagePattern string = "startingIndex=[0-9]*"
 	// Agriculture string = "WSJMXUSAGRI"
 	// FinancialServices string = "WSJMXUSFCL"
 )

  var Industries = map[string]string{
 		"Agricultre": "WSJMXUSAGRI",
 		"Automotive": "WSJMXUSAUTO",
 		"Basic Materials/Resources": "WSJMXUSBSC",
 		"Business/Consumer Services": "WSJMXUSCYC",
 		"Consumer Goods": "WSJMXUSNCY",
 		"Energy": "WSJMXUSENE",
 		"Financial Services": "WSJMXUSFCL",
 		"Health Care/Life Sciences": "WSJMXUSHCR",
 		"Industrial Goods": "WSJMXUSIDU",
 		"Leisure/Arts/Hospitality": "WSJMXUSLEAH",
 		"Media/Entertainment": "WSJMXUSMENT",
 		"Real Estate/Construction": "WSJMXUSRECN",
 		"Retail/Wholesale": "WSJMXUSRTWS",
 		"Technology": "WSJMXUSTEC",
 		"Telecommunication Services": "WSJMXUSTEL",
 		"Transportation/Logistics": "WSJMXUSTRSH",
 		"Utilities": "WSJMXUSUTI",
 	}

func getStockSymbolsByIndustry(industry string, wg *sync.WaitGroup) {
	start := time.Now()
	fmt.Printf("%s starting now ...\n", industry)
	defer func() {
		elapsed := time.Since(start)
 		fmt.Printf("Industry: %s finished in: %s\n", industry, elapsed)
 		wg.Done()
	}()

	// Decide the number of concurrent processes
	urlMain := fmt.Sprintf(Url, industry, 0)
	re := regexp.MustCompile(PagePattern)
	resp, err := http.Get(urlMain)
	if err != nil {
		// error
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//error
	}
	res := re.FindAll(b, -1)
	/* Long expression, its just getting the last page number of the set of Industry Pages
	from here we work backwards or forwards in the For loop depending if there is any content */
	lastPage, err := strconv.Atoi(strings.Split(string(res[len(res)-1]), "=")[1])
	// lastPage = lastPage
	oMidPage := math.Ceil(float64(lastPage) / 50 / 2) * 50
	re = regexp.MustCompile(StockPattern)

	// oMidPage := midPage
	midPage := float64(lastPage)
	// fmt.Printf("LAST PAGE: %g\n", midPage)
	run := true
	// counter := 0
	pagesSeen := map[int]bool{}
	for run {
		count, ok := pagesSeen[int(oMidPage)]
		if count && ok {
			// fmt.Printf("Run truth: %g Seen: %v\n", oMidPage, pagesSeen)
			// oMidPage -= 50
			break
		}
		// counter++
		url := fmt.Sprintf(Url, industry, int(oMidPage))
		
		// fmt.Printf("URL read : %s Truth value: %t\n", url, run)
		resp, err := http.Get(url)
		if err != nil {
			
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//error
		}

		res := re.FindAll(b, 1)
		// The case where you have to go higher
		if len(res) > 0 {
			diff := math.Ceil((math.Abs(midPage - oMidPage) / 50 / 2)) * 50
			pagesSeen[int(oMidPage)] = true
			// fmt.Printf("higher: %g Seen: %v\n", diff, pagesSeen)
			midPage = oMidPage
			oMidPage += diff
			
		} 
		// The case where you have to go lower
		if len(res) == 0 {
			diff := math.Ceil((math.Abs(midPage - oMidPage) / 50 / 2)) * 50 
			pagesSeen[int(oMidPage)] = false
			// fmt.Printf("lower: %g Seen: %v\n", diff, pagesSeen)
			midPage = oMidPage
			oMidPage -= diff
			// pagesSeen[int(oMidPage)]++
		}
	}
	// fmt.Printf("Counter %d TRUE value: %t\n", counter, run)
	concurrency := int(oMidPage / 50) + 1
	var wgI sync.WaitGroup
	wgI.Add(concurrency)
	// fmt.Printf("Concurrency %s is %d\n", industry, concurrency)
 	for i := 0; i < concurrency; i++ {
 		pageNo := i * 50
 		// fmt.Printf("go routine %s pages: %d\n", industry, pageNo)
 		go getStockSymbolsByPage(industry, pageNo, &wgI)
 	}

	wgI.Wait()
}

func getStockSymbolsByPage(industry string, page int, wgI *sync.WaitGroup) {
	// start := time.Now()
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
 	re := regexp.MustCompile(StockPattern)
	for _, res := range re.FindAll(b, -1) {
 		s := fmt.Sprintf("%s", res)
 		i, j := strings.Index(s, ">") + 1, strings.Index(s, "</td>")
 		x, y := strings.Index(s, "<div>") + 5, strings.Index(s, "</div>")
 		Symbols = append(Symbols,StockSymbol{Code: s[i:j], Description: s[x:y]})
	}
 	resp.Body.Close()

 	data, err := json.Marshal(Symbols)
 	if err != nil {
 		fmt.Println("Could not marshal struct to JSON")
 	}
 	fileName := "/tmp/" + industry + "_" + strconv.Itoa(page) + ".json"
 	if err := ioutil.WriteFile(fileName, data, os.ModePerm); err != nil {
 		fmt.Printf("Error writing file: %s with error: %s\n", fileName, err)
 	}
 	defer wgI.Done()
}

func main() {
 	mainStart := time.Now()
	
	defer func() {
		mainElapsed := time.Since(mainStart)
		fmt.Printf("\nTotal time: %s for %d number of Industries\n", mainElapsed, len(Industries))
	}()

 	var wg sync.WaitGroup
 	wg.Add(len(Industries))

 	for _, industry := range Industries {
 		go getStockSymbolsByIndustry(industry, &wg)
 	}

 	wg.Wait()
 }