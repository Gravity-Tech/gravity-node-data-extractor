// Package classification Gravity Extractor RPC API.
//
// This application represents viable extractor methods.
// Declared methods are compulsory for appropriate extractor functioning.
//
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: extractor.gravityhub.org
//     BasePath: /
//     Version: 1.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: venlab.dev <shamil@venlab.dev> https://venlab.dev
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - api_key:
//
//     SecurityDefinitions:
//     api_key:
//          type: apiKey
//          name: KEY
//          in: header
//
//     Extensions:
//     x-meta-value: value
//     x-meta-array:
//       - value1
//       - value2
//     x-meta-array-obj:
//       - name: obj
//         value: field
//
// swagger:meta
package main

import (
	"flag"
	"fmt"
	c "github.com/Gravity-Tech/gravity-node-data-extractor/v2/controller"
	r "github.com/Gravity-Tech/gravity-node-data-extractor/v2/router"
	"net/http"
)
var port, extractorTag, symbolPair, apiKey, extractorType string

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func init() {
	flag.StringVar(&port, "port", "8090", "Port to run on")
	flag.StringVar(&extractorTag, "tag", "latest", "Extractor version tag")
	flag.StringVar(&symbolPair, "pair", "WAVESBTC", "Pair symbol appropriate to Binance API")
	flag.StringVar(&apiKey, "api", "NONE", "Binance API Key")
	flag.StringVar(&extractorType, "type", "binance", "Extractor Type")

	flag.Parse()
}

func main () {
	tagController := &c.ParamsController{ Tag: extractorTag, SymbolPair: symbolPair, ApiKey: apiKey, ExtractorType: extractorType }
	respController := &c.ResponseController{ TagDelegate: tagController }

	http.HandleFunc(r.GetExtractedData, respController.GetExtractedData)
	http.HandleFunc(r.GetExtractorInfo, respController.GetExtractorInfo)
	http.HandleFunc(r.GetAggregated, respController.Aggregate)

	fmt.Printf("Listening on port: %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}