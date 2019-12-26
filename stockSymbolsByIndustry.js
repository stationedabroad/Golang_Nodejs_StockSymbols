"use strict"

const request = require('request-promise')
const matchAll = require('match-all')

const URL = "http://bigcharts.marketwatch.com/industry/bigcharts-com/stocklist.asp?"
const FinancialServices = "WSJMXUSFCL"
const Agriculture = "WSJMXUSAGRI"
const tickerUrl = URL + "Symb=" + FinancialServices + "&startingIndex="
const re = "<td class=\"symb-col\">[A-Za-z0-9.]*</td>\\s*<td class=\"name-col\"><div>.*</div>";

function readIndustryTicker(url) {
	let startRead = new Date()
	return request.get(url)
		.then(response => {
			let result = [...response.matchAll(re)];
			// console.log("Done processing: ", url)
			let res = result.map(entry => {
				var startIdxCode = entry[0].indexOf(">") + 1
				var endIdxCode = entry[0].indexOf("</td>")
				var startIdxDesc = entry[0].indexOf("<div>") + 5
				var endIdxDesc = entry[0].indexOf("</div>")
				return {
					"symbol_code": entry[0].substring(startIdxCode, endIdxCode),
					"company_desc": entry[0].substring(startIdxDesc, endIdxDesc)
				}
			});
			console.info("Execution of: %s\t entries: %d\t time elapsed: %dms", url, res.length, new Date() - startRead)
			return res
		});
};

let start = new Date()
var stockFuncs = new Array()
for(let i = 0; i <= 4400; i += 50) {
	var pageTickerUrl =  tickerUrl + i
	stockFuncs.push(readIndustryTicker(pageTickerUrl))
}

console.log("Sent jobs ...")

var promiseAllStock = Promise.all(stockFuncs)
promiseAllStock.then(function(results) {
	console.info("Total Execution time: %dms", new Date() - start)
	results.forEach(function(tab) {
		console.info("\n", tab.length)
	});
});

