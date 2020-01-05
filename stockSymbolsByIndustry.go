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

	"runtime/pprof"
	"flag"
	"log"
	"runtime"
)

type StockSymbol struct{
	Code string `json:"symbol_code"`
	Description string `json:"company_name"`
}

const (
 	Url string = "http://bigcharts.marketwatch.com/industry/bigcharts-com/stocklist.asp?Symb=%s&startingIndex=%d"
 	StockPattern string = "<td class=\"symb-col\">[A-Za-z0-9.]*</td>\\s*<td class=\"name-col\"><div>.*</div>"
 	PagePattern string = "startingIndex=[0-9]*"
 )

var Industries = map[string]string{
 		"Agriculture": "WSJMXUSAGRI",
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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func getStockSymbolsByIndustry(industry string, wg *sync.WaitGroup) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
 		fmt.Printf("Industry: %s\t finished in: %s\n", industry, elapsed)
 		wg.Done()
	}()

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
	lastPage, err := strconv.Atoi(strings.Split(string(res[len(res)-1]), "=")[1])
	nextPage := math.Ceil(float64(lastPage) / 50 / 2) * 50
	prevPage := float64(lastPage)
	re = regexp.MustCompile(StockPattern)
	run := true
	pagesSeen := map[int]bool{}
	for run {
		count, ok := pagesSeen[int(nextPage)]
		if count && ok {
			break
		}

		url := fmt.Sprintf(Url, industry, int(nextPage))
		resp, err := http.Get(url)
		if err != nil {
			// error
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//error
		}

		res := re.FindAll(b, 1)
		if len(res) > 0 {
			diff := math.Ceil((math.Abs(prevPage - nextPage) / 50 / 2)) * 50
			pagesSeen[int(nextPage)] = true
			prevPage = nextPage
			nextPage += diff
			
		} 
		if len(res) == 0 {
			diff := math.Ceil((math.Abs(prevPage - nextPage) / 50 / 2)) * 50 
			pagesSeen[int(nextPage)] = false
			prevPage = nextPage
			nextPage -= diff
		}
	}
	concurrency := int(nextPage / 50) + 1
	var wgI sync.WaitGroup
	wgI.Add(concurrency)
 	for i := 0; i < concurrency; i++ {
 		pageNo := i * 50
 		go getStockSymbolsByPage(industry, pageNo, &wgI)
 	}

	wgI.Wait()
}

func getStockSymbolsByPage(industry string, page int, wgI *sync.WaitGroup) {
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
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could bot create cpu profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
		}
	}
	
	defer pprof.StopCPUProfile()

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

 	if *memprofile != "" {
 		f, err := os.Create(*memprofile)
 		if err != nil {
 			log.Fatal("could not create memory profile: ", err)
 		}
 		defer f.Close()
 		runtime.GC()
 		if err := pprof.WriteHeapProfile(f); err != nil {
 			log.Fatal("could not write memory profile: ", err)
 		}
 	}
 }