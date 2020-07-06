package model

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"fmt"
)

// swagger:model
type BinancePriceIndexResponse struct {
	Symbol string `json:"symbol"`
	Price string `json:"price"`
	CalcTime int64 `json:"calcTime"`
}

// swagger:model
type BinancePriceExtractor struct {
	Tag, SymbolPair, ApiKey string
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
func (e *BinancePriceExtractor) Data() (interface{}, interface{}) {
	return e.Price()
}
func (e *BinancePriceExtractor) endpoint() string {
	return "https://api.binance.com/sapi/v1/margin/priceIndex"
}
func (e *BinancePriceExtractor) headers () http.Header {
	dict := make(http.Header)
	dict["X-MBX-APIKEY"] = []string { e.ApiKey }
	return dict
}

func (e *BinancePriceExtractor) requestPrice() *BinancePriceIndexResponse {
	headers := e.headers()
	endpoint := e.endpoint()
	endpoint = fmt.Sprintf("%v?symbol=WAVESBTC", endpoint)

	url, _ := url.ParseRequestURI(endpoint)

	request := http.Request{
		Method:           "GET",
		URL:              url,
		Header:           headers,
	}

	resp, err := http.DefaultClient.Do(&request)

	defer func () {
		_ = resp.Body.Close()
	}()

	if err != nil {
		fmt.Printf("Error occured: %v \n", err)
		return nil
	}

	var result BinancePriceIndexResponse

	byteValue, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(byteValue, &result)

	return &result
}
func (e *BinancePriceExtractor) encodeFloat (buf []RawData, f float64) {
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
}
func (e *BinancePriceExtractor) decodeFloat (buf *[]RawData) float64 {
	extr := binary.BigEndian.Uint64(*buf)
	fl := math.Float64frombits(extr)
	return fl
}
func (e *BinancePriceExtractor) extractData(params interface{}) []RawData {
	extracted := make([]RawData, 8)
	castedParams := params.(*BinancePriceIndexResponse)

	floatCurrentPrice, err := strconv.ParseFloat(castedParams.Price, 64)

	if err != nil {
		fmt.Printf("Failed to parse to float: %v \n", err)
		return extracted
	}
	e.encodeFloat(extracted, floatCurrentPrice)

	fmt.Printf("Raw: %v; Price: %v; Uint: %v \n", extracted, floatCurrentPrice, math.Float64bits(floatCurrentPrice))
	return extracted
}
func (e *BinancePriceExtractor) mapData(data []RawData) interface{} {
	return e.decodeFloat(&data)
}
func (e *BinancePriceExtractor) Info() *ExtractorInfo {
	return &ExtractorInfo{ DataFeedTag: e.DataFeedTag(), Description: e.Description() }
}