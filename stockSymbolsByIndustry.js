"use strict"

const request = require('request-promise')
const matchAll = require('match-all')

const URL = "http://bigcharts.marketwatch.com/industry/bigcharts-com/stocklist.asp?"
const Pattern = "<td class=\"symb-col\">[A-Za-z0-9.]*</td>\\s*<td class=\"name-col\"><div>.*</div>"
const FinancialServices = "WSJMXUSFCL"
const Agriculture = "WSJMXUSAGRI"
const regex = new RegExp("<td class=\"symb-col\">[A-Za-z0-9.]*</td>\\s*<td class=\"name-col\"><div>.*</div>")
let re = "<td class=\"symb-col\">[A-Za-z0-9.]*</td>\\s*<td class=\"name-col\"><div>.*</div>";

function readIndustryTicker(url) {
	return request.get(url)
		.then(response => {
			// console.log(response)
			let result = [...response.matchAll(re)];
			// let results = new Array();
			// var res;
			// do {
			// 	res = regex.exec(response)
			// 	if (res) {
			// 		// results.push(res[0])
			// 		console.log(res[0])
			// 	}
			// } while (res)
			return result.map(entry => entry[0]);
		});
};

let tickerUrl = URL + "Symb=" + Agriculture + "&startingIndex=" + 0
console.log(readIndustryTicker(tickerUrl)
			.then(data => {
				console.log(data);
			}));
// console.log(tickerUrl)