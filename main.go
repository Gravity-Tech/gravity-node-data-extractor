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

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/config"

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

var isDirect bool
var port, extractorType, configName string

func init() {
	flag.StringVar(&port, "port", "8090", "Port to run on")
	flag.BoolVar(&isDirect, "direct", false, "Is direct swap extractor")
	flag.StringVar(&extractorType, "type", string(WavesSource), "Extractor Type")
	flag.StringVar(&configName, "config", config.MainConfigFile, "Config file name")

	flag.Parse()
}

func main() {
	ctx := context.Background()
	var extractor extractors.IExtractor
	var err error

	cfg, err := config.ParseMainConfig(configName)

	if err != nil {
		panic(err)
	}

	println(extractorType)
	switch ExtractorType(extractorType) {
	case BinanceWavesBtc:
		extractor = &binance.Extractor{}
	case WavesSource:
		extractor, err = susy.New(
			cfg.SourceNodeURL,
			cfg.DestinationNodeURL,
			cfg.LUPortAddress,
			cfg.IBPortAddress,
			cfg.SourceDecimals,
			cfg.DestinationDecimals,
			ctx,
			susy.RequestDirection { Kind: susy.WavesSource, IsDirect: isDirect },
		)
	case EthereumSource:
		extractor, err = susy.New(
			cfg.SourceNodeURL,
			cfg.DestinationNodeURL,
			cfg.LUPortAddress,
			cfg.IBPortAddress,
			cfg.SourceDecimals,
			cfg.DestinationDecimals,
			ctx,
			susy.RequestDirection { Kind: susy.WavesSource, IsDirect: isDirect },
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
