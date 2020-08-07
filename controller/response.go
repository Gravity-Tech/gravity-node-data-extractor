package controller

import (
	"encoding/json"
	"fmt"
	m "github.com/Gravity-Tech/gravity-node-data-extractor/v2/model"
	"net/http"
)

type ResponseController struct {
	TagDelegate *ParamsController
}

func (rc *ResponseController) extractorEnumerator () *m.ExtractorEnumerator {
	return m.DefaultExtractorEnumerator
}

func (rc *ResponseController) aggregator() m.Aggregator {
	return &m.BinanceAggregator{}
}

func (rc *ResponseController) extractor () *m.ExtractorProvider {
	enumerator := rc.extractorEnumerator()

	var extractor m.IExtractor

	switch rc.TagDelegate.ExtractorType {
	case enumerator.Metal:
		extractor = &m.MetalCurrencyMetalExtractor{
			Tag:        rc.TagDelegate.Tag,
			MetalIndex: "XAU",
		}
	case enumerator.Binance:
		fallthrough
	default:
		extractor = &m.BinancePriceExtractor{ Tag: rc.TagDelegate.Tag, SymbolPair: rc.TagDelegate.SymbolPair, ApiKey: rc.TagDelegate.ApiKey }
	}

	return &m.ExtractorProvider{ Current: extractor }
}

func addBaseHeaders (headers http.Header) {
	headers.Add("Content-Type", "application/json")
}
// swagger:route GET /extracted Extractor getExtractedData
//
// Extracts mapped data
//
// No additional info
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Deprecated: false
//
//     Security:
//       api_key:
//
//     Responses:
//       200: BinancePriceIndexResponse
func (rc *ResponseController) GetExtractedData (w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" { return }

	extractor := rc.extractor().Current

	_, extractedData := extractor.Data()

	addBaseHeaders(w.Header())

	bytes, _ := json.Marshal(&extractedData)

	_, _ = fmt.Fprint(w, string(bytes))
}

// swagger:route GET /raw Extractor getRawData
//
// Resolves raw data
//
// No additional info
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Deprecated: false
//
//     Security:
//       api_key:
//
//     Responses:
//       200: RawData
func (rc *ResponseController) GetRawData (w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" { return }

	extractor := rc.extractor().Current

	rawResponse, _ := extractor.Data()

	addBaseHeaders(w.Header())

	bytes, _ := json.Marshal(&rawResponse)

	_, _ = fmt.Fprint(w, string(bytes))
}


// swagger:route GET /info Extractor getExtractorInfo
//
// Returns extractor common info
//
// No additional info
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Deprecated: false
//
//     Security:
//       api_key:
//
//     Responses:
//       200: ExtractorInfo
func (rc *ResponseController) GetExtractorInfo (w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" { return }

	extractor := rc.extractor().Current
	extractorInfo := extractor.Info()

	addBaseHeaders(w.Header())

	bytes, _ := json.Marshal(&extractorInfo)

	_, _ = fmt.Fprint(w, string(bytes))
}


func (rc *ResponseController) Aggregate (w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" { return }

	type requestType = string
	const (
		typeInt64 requestType = "int64"
		typeFloat64 = "float64"
		typeString = "string"
	)

	type requestBody struct {
		Type string `json:"type"`
		Values []interface{} `json:"values"`
	}
	var paramsBody requestBody

	decoder := json.NewDecoder(req.Body)
	aggregator := rc.aggregator()
	var result interface{}

	addBaseHeaders(w.Header())

	if err := decoder.Decode(&paramsBody); err != nil {
		_, _ = fmt.Fprint(w, fmt.Errorf("Invalid body", err))

		return
	}

	switch paramsBody.Type {
	case typeInt64:
		result = aggregator.AggregateInt(paramsBody.Values)
		break
	case typeFloat64:
		result = aggregator.AggregateFloat(paramsBody.Values)
		break
	case typeString:
		result = aggregator.AggregateString(paramsBody.Values)
		break
	}

	_, _ = fmt.Fprint(w, result)

}
