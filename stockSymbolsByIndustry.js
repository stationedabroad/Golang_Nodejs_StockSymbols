"use strict"

const request = require('request-promise')
const http = require('http')
const matchAll = require('match-all')
const fs = require('fs')

const Url = "http://bigcharts.marketwatch.com/industry/bigcharts-com/stocklist.asp?Symb="
const StockPattern = "<td class=\"symb-col\">[A-Za-z0-9.]*</td>\\s*<td class=\"name-col\"><div>.*</div>"
const PagePattern = "startingIndex=[0-9]*"
const writeDirectory = "/home/sulman/Downloads/node-v12.14.0-linux-x64/bin/tmp/"
const Industries = {
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

function getSyncRequest(url) {
	let chunks_of_data = [];
	return new Promise((res, rej) => {
		http.get(url, response => {
			response.on('data', (chunk) => {
				chunks_of_data.push(chunk);
			})

			response.on('end', () => {
				let resp_body = Buffer.concat(chunks_of_data);
				res(resp_body.toString());
			});

			response.on('error', (error) => {
				rej(error);
			});
		});
	});
}

function readTickerByIndustry(industry, url) {
	let industryStart = new Date()
	return new Promise((resolve, reject) => {
	request.get(url + 0)
		.then(async function(response) {
			let mainPageResults = [...response.matchAll(PagePattern)];
			let lastPage = mainPageResults[mainPageResults.length-1][0].split('=')[1]
			let nextPage = Math.ceil(lastPage / 50 / 2) * 50
			let prevPage = lastPage

			let pagesSeen = new Map()
			let run = true
			let pageUrl = url + nextPage
			var diff = 0
			while(run) {
				if(pagesSeen.has(nextPage) && pagesSeen.get(nextPage)) {
					break
				}
				var httpPromise = getSyncRequest(url + nextPage);
				var httpResponse = await httpPromise;

				if([...httpResponse.matchAll(StockPattern)].length > 0) {
					diff = Math.ceil((Math.abs(prevPage - nextPage) / 50 / 2)) * 50
					pagesSeen.set(nextPage, true)
					prevPage = nextPage
					nextPage += diff
				} else {
					diff = Math.ceil((Math.abs(prevPage - nextPage) / 50 / 2)) * 50
					pagesSeen.set(nextPage, false)
					prevPage = nextPage
					nextPage -= diff				
				}					
			}

			let pageFuncs = new Array();	
			let concurrency = nextPage / 50 + 1
			for(let i = 0; i < concurrency; i++) {
				var page = i * 50
				pageFuncs.push(readTickerByPage(industry, url, page)) 
			}

			let promiseAllPages = Promise.all(pageFuncs)
			promiseAllPages.then((pageResults) => {
				console.info("Industry: %s\t Finished in: %dms",  industry, new Date() - industryStart)
				resolve();
			});
		});
	});
};

function readTickerByPage(industry, url, page) {
	let startRead = new Date()
	return new Promise((resolve, reject) => { 
	request.get(url + page)
		.then(response => {
			let result = [...response.matchAll(StockPattern)];
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
			let resJson = JSON.stringify(res, null, 4);
			fs.writeFile(writeDirectory + industry + "_" + page + ".json", resJson, (err) => {
				if(err) {
					console.log("Could not write: %s with page: %d", industry, page )
					reject();
				}
			});
		resolve();
		});
	});
};

let hrstart = process.hrtime()
var stockFuncs = new Array()
for(const industry in Industries) {
	var pageTickerUrl =  Url + Industries[industry] + "&startingIndex="
	stockFuncs.push(readTickerByIndustry(Industries[industry], pageTickerUrl))
}

console.log("Sent %d jobs ...", stockFuncs.length)

var promiseAllStock = Promise.all(stockFuncs)
promiseAllStock.then(function(results) {
	let hrend = process.hrtime(hrstart)
	console.info("\nTotal execution time: %ds %dms", hrend[0], hrend[1] / 1000000)
});


