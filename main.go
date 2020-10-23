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
	EthereumSource  ExtractorType = "ethereum-source"
)

type ExtractorType string

var port, extractorType string

func init() {
	flag.StringVar(&port, "port", "8090", "Port to run on")
	flag.StringVar(&extractorType, "type", "waves-source", "Extractor Type")

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
			"http://app-cccba780-bc8c-4cfb-9cba-d566055a5e2c.cls-dec3c32b-4f06-462f-b827-dee931d39a72.ankr.com",
			"3Mk3SUp8Cs7f6958reTnXsmr1QTJ26wcuRN",
			"0x042AF5c312489Be04882c1ADBfA0AD9E68d2e66B",
			ctx,
			susy.WavesSourceLock,
		)
	case EthereumSource:
		extractor, err = susy.New(
			"https://nodes-stagenet.wavesnodes.com",
			"https://ropsten.infura.io/v3/663ad61d27254aac874ba7fc298e0956",
			"3MnZnDHsDTRvsrFrVtUp5zVuWybqHpHFMsy",
			"0x2f40ac805df8fa8c862d16412a92c439b4a90675",
			ctx,
			susy.EthereumSourceBurn,
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
