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
	"errors"
	"flag"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/config"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/binance"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/server"
)

var port, configName string
var extractorType extractors.ExtractorType

func init() {
	flag.StringVar(&port, "port", "8090", "Port to run on")
	flag.StringVar((*string)(&extractorType), "type", string(susy.WavesToEthDirect), "Extractor Type")
	flag.StringVar(&configName, "config", config.MainConfigFile, "Config file name")

	flag.Parse()
}

func main() {
	var extractor extractors.IExtractor
	var err error

	cfg, err := config.ParseMainConfig(configName)

	if err != nil {
		panic(err)
	}

	println(extractorType)
	switch extractorType {
	case binance.BinanceWavesBtc:
		extractor = &binance.Extractor{}
	case susy.WavesToEthDirect, susy.WavesToEthReverse, susy.EthToWavesDirect, susy.EthToWavesReverse:
		extractor, err = susy.New(
			cfg.SourceNodeURL,
			cfg.DestinationNodeURL,
			cfg.LUPortAddress,
			cfg.IBPortAddress,
			cfg.SourceDecimals,
			cfg.DestinationDecimals,
			extractorType,
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
