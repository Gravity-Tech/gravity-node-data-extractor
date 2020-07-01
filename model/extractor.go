package model

import (
	"fmt"
	"strconv"
	"time"
)

// swagger:model
type RawData = byte

type IExtractor interface {
	DataFeedTag() string
	Description() string
	extractData(params interface{}) []RawData
	mapData(extractedData []RawData) interface{}
}

// swagger:model
type ExtractorInfo struct {
	Description string
	DataFeedTag string
}

// swagger:model
type BinancePriceIndexResponse struct {
	Symbol string `json:"symbol"`
	Price string `json:"price"`
	CalcTime int64 `json:"calcTime"`
}

// swagger:model
type BinancePriceExtractor struct {
	Tag, SymbolPair string
}

func (e *BinancePriceExtractor) DataFeedTag() string {
	return fmt.Sprintf("binance-WAVES_BTC:%v", e.Tag)
}
func (e *BinancePriceExtractor) Description() string {
	return "This extractor resolves price data for WAVES_BTC pair presented in decimal"
}
func (e *BinancePriceExtractor) Price() (*BinancePriceIndexResponse, float64) {
	priceResponse := e.requestPrice()
	extractedPrice := e.extractData(priceResponse)
	mappedData := e.mapData(extractedPrice).(float64)

	return priceResponse, mappedData
}

func (e *BinancePriceExtractor) requestPrice() *BinancePriceIndexResponse {
	return &BinancePriceIndexResponse{
		Symbol:   e.SymbolPair,
		Price:    "0.05",
		CalcTime: time.Time{}.Unix(),
	}
}
func (e *BinancePriceExtractor) extractData(params interface{}) []RawData {
	extracted := make([]RawData, 1)
	castedParams := params.(BinancePriceIndexResponse)

	floatCurrentPrice, err := strconv.ParseFloat(castedParams.Price, 64)

	if err != nil {
		extracted[0] = RawData(floatCurrentPrice)
	}

	return extracted
}
func (e *BinancePriceExtractor) mapData(data []RawData) interface{} {
	extractedPrice := data[0]
	return float64(extractedPrice)
}
func (e *BinancePriceExtractor) Info() *ExtractorInfo {
	return &ExtractorInfo{ DataFeedTag: e.DataFeedTag(), Description: e.Description() }
}