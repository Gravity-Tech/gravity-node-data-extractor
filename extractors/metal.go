package extractors

/*
import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"time"
)

// swagger:model
type MetalCurrencyHistoryRate struct {
	//BankBuyAt": 2917.5,
	//"BankSellAt": 3393,
	//"CbRate": 3031.25,
	//"OnDate": "2020-01-06T00:00:00+00:00"
	BankBuyAt float64
	BankSellAt float64
	CbRate float64
	OnDate string
}

var dateFormat = "mm.dd.yyyy"

var defaultContextItemId = "%7BEAAE3BDA-B447-4642-8561-DBF5C8E28AFA%7D"

// swagger:model
type MetalCurrencyExchangeResponse struct {
	MetalHistoryRate *[]MetalCurrencyHistoryRate
	DifferenceRubles float64
	DifferencePercents float64
}

type MetalCurrencyMetalExtractor struct {
	Tag, MetalIndex string
}
func (e *MetalCurrencyMetalExtractor) endpoint() string {
	return "https://www.vtb.ru/api/currency-exchange/metal-calculator-info"
}
func (e *MetalCurrencyMetalExtractor) DataFeedTag() string {
	return fmt.Sprintf("metal-GOLD:%v", e.Tag)
}
func (e *MetalCurrencyMetalExtractor) Description() string {
	return "This extractor resolves metal price data for specific date range"
}
func (e *MetalCurrencyMetalExtractor) RequestMetalIndex() *MetalCurrencyExchangeResponse {
	return e.requestData(e.MetalIndex, 100)
}
func (e *MetalCurrencyMetalExtractor) formatDate(t time.Time) string {
	return fmt.Sprintf("%02d.%02d.%d", t.Month()-1, t.Day(), t.Year())
}
func (e *MetalCurrencyMetalExtractor) requestData(metalCode string, amount int64) *MetalCurrencyExchangeResponse {
	endpoint := e.endpoint()

	respUrl, urlErr := url.Parse(endpoint)

	if urlErr != nil {
		fmt.Printf("Error occured on parse: %v\n", urlErr)
		return nil
	}
	if metalCode == "" {
		metalCode = "XAU"
	}
	if amount < 1 {
		amount = 100
	}

	oneWeek := time.Hour * 24 * 7
	dateFrom := time.Now().Add(-oneWeek + time.Hour * 24)
	dateTo := dateFrom.Add(oneWeek - time.Hour * 24)
	var queryParams = map[string]string {
		"amount": 			 fmt.Sprintf("%v", amount),
		"contextItemId":     defaultContextItemId,
		"dateFrom":          e.formatDate(dateFrom),
		"dateTo":            e.formatDate(dateTo),
		"isRublesInput":     fmt.Sprintf("%v", true),
		"selectedMetalCode": metalCode,
	}

	params := make(url.Values)
	for key, value := range queryParams {
		params.Set(key, value)
	}

	provided := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	fmt.Printf("PRO: %v", provided)
	respUrl, _ = url.Parse(provided)
	fmt.Printf("PARAMS: %v", respUrl.Query())

	request := http.Request{
		Method:           "GET",
		URL:              respUrl,
		Header: map[string][]string {
			"cache-control": []string { "no-cache" },
		},
	}

	resp, err := http.DefaultClient.Do(&request)

	defer func () {
		_ = resp.Body.Close()
	}()

	if err != nil {
		fmt.Printf("Error occured: %v \n", err)
		return nil
	}

	var result MetalCurrencyExchangeResponse

	byteValue, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(byteValue, &result)

	return &result
}
func (e *MetalCurrencyMetalExtractor) Data() (interface{}, interface{}) {
	raw := e.RequestMetalIndex()
	return raw, raw
}
func (e *MetalCurrencyMetalExtractor) extractData(data interface{}) []RawData {
	castedParams := data.(*MetalCurrencyExchangeResponse)
	res, _ := json.Marshal(castedParams)
	return res
}
func (e *MetalCurrencyMetalExtractor) encodeFloat (buf []RawData, f float64) {
	binary.BigEndian.PutUint64(buf, math.Float64bits(f))
}
func (e *MetalCurrencyMetalExtractor) decodeFloat (buf *[]RawData) float64 {
	extr := binary.BigEndian.Uint64(*buf)
	fl := math.Float64frombits(extr)
	return fl
}
func (e *MetalCurrencyMetalExtractor) mapData(data []RawData) interface{} {
	return e.decodeFloat(&data)
}
func (e *MetalCurrencyMetalExtractor) Info() *ExtractorInfo {
	return &ExtractorInfo{ DataFeedTag: e.DataFeedTag(), Description: e.Description() }
}*/
