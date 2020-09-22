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
	"context"
	"errors"
	"flag"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/binance"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/server"
)

const (
	BinanceWavesBtc ExtractorType = "binance-waves-btc"
	WavesSource     ExtractorType = "waves-source"
)

type ExtractorType string

var port, extractorType string

func init() {
	flag.StringVar(&port, "8port", "8090", "Port to run on")
	flag.StringVar(&extractorType, "type", "binance-waves-btc", "Extractor Type")

	flag.Parse()
}

func main() {
	ctx := context.Background()
	var extractor extractors.IExtractor
	var err error
	switch ExtractorType(extractorType) {
	case BinanceWavesBtc:
		extractor = &binance.Extractor{}
	case WavesSource:
		extractor, err = susy.New(
			"https://nodes-stagenet.wavesnodes.com",
			"https://ropsten.infura.io/v3/663ad61d27254aac874ba7fc298e0956",
			"3MdQFC6chdxJ2WrxYV4ZidmutZdpzea1Kqp",
			"0xbd54863045214cFc03385305758407EB27c484C0",
			ctx,
		)
	default:
		panic(errors.New("invalid "))
	}

	if err != nil {
		panic(err)
	}

	server := server.New(extractor)
	err = server.Start(port)
	if err != nil {
		panic(err)
	}
}
