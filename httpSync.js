"use strict"

const http = require('http')
const request = require('request-promise')

// var chunks_of_data = [];

function getRequest(url) {
	let chunks_of_data = [];
	return new Promise((res, rej) => {
		http.get(url, response => {
			response.on('data', (chunk) => {
				console.log("Called on chunk receieve ...", url)
				chunks_of_data.push(chunk);
			})

			response.on('end', () => {
				console.log("Called on End ...")
				let resp_body = Buffer.concat(chunks_of_data);
				res(resp_body.toString());
			});

			response.on('error', (error) => {
				rej(error);
			});
		});
	});
}

async function caller() {
	let urls = ["http://bigcharts.marketwatch.com/industry/bigcharts-com/stocklist.asp?Symb=WSJMXUSBSC&startingIndex=0",
				"http://bigcharts.marketwatch.com/industry/bigcharts-com/stocklist.asp?Symb=WSJMXUSTEC&startingIndex=50"]
	try {
		for(let i = 0; i < urls.length; i++) {
			console.log("Calling url: ", urls[i])
			let http_promise = getRequest(urls[i]);
			let http_resp = await http_promise;

			console.log("Called after wait ...")
			console.log(http_resp.slice(0,500));
			console.log("\n\n\n")
		}
	} catch (error) {
		console.error(error)
	}
}

console.log(caller())


