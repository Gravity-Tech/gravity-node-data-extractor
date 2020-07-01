package controller

import (
	m "../model"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseController struct {
	TagDelegate *ParamsController
}

func (rc *ResponseController) extractor () m.IExtractor {
	return &m.BinancePriceExtractor{ Tag: rc.TagDelegate.Tag, SymbolPair: rc.TagDelegate.SymbolPair }
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
	extractor := rc.extractor().(*m.BinancePriceExtractor)
	_, extractedData := extractor.Price()

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
	extractor := rc.extractor().(*m.BinancePriceExtractor)
	rawResponse, _ := extractor.Price()

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
	extractor := rc.extractor().(*m.BinancePriceExtractor)
	extractorInfo := extractor.Info()

	addBaseHeaders(w.Header())

	bytes, _ := json.Marshal(&extractorInfo)

	_, _ = fmt.Fprint(w, string(bytes))
}
