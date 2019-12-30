package main

import (
	"fmt"
	"regexp"
	"net/http"
	"io/ioutil"
	"strings"
	"strconv"
	"math"
)

const (
 	Url string = "http://bigcharts.marketwatch.com/industry/bigcharts-com/stocklist.asp?Symb=%s&startingIndex=%d"
 	// StockPattern string = "<td class=\"symb-col\">[A-Za-z0-9.]*</td>\\s*<td class=\"name-col\"><div>.*</div>"
 	PagePattern string = "startingIndex=[0-9]*"
 	Agriculture string = "WSJMXUSAGRI"
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

func main() {
	urlMain := fmt.Sprintf(Url, Agriculture, 0)
	re := regexp.MustCompile(PagePattern)
	resp, err := http.Get(urlMain)
	if err != nil {
		// error
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//error
	}
	// for _, r := range re.FindAll(b, -1) {
	// 	num, err := strconv.Atoi(strings.Split(string(r), "=")[1])
	// 	if err == nil {
	// 		fmt.Println(num/2)
	// 	}
	// }
	sl := []int{}
	sl = append(sl, 200)
	sl = append(sl, 400)
	res := re.FindAll(b, -1)
	num, err := strconv.Atoi(strings.Split(string(res[len(res)-1]), "=")[1])
	page := math.Ceil(float64(num)/50/2) * 50
	// page = math.Ceil(page) * 50
	fmt.Println(num, page, sl)
	fmt.Println(len(Industries))
}