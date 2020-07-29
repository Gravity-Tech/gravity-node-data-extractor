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

func (rc *ResponseController) extractorEnumerator () *m.ExtractorEnumerator {
	return m.DefaultExtractorEnumerator
}

func (rc *ResponseController) extractor () *m.ExtractorProvider {
	enumerator := rc.extractorEnumerator()

	var extractor m.IExtractor

	switch rc.TagDelegate.ExtractorType {
	case enumerator.WavesChain:
		extractor = &m.WavesChainExtractor{}
	case enumerator.EthereumChain:
		extractor = &m.EthChainExtractor{}
	default:
		extractor = &m.WavesChainExtractor{}
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
	extractor := rc.extractor().Current

	_, extractedData := extractor.Data()

	addBaseHeaders(w.Header())

	bytes, _ := json.Marshal(&extractedData)

	_, _ = fmt.Fprint(w, string(bytes))
}


func (rc *ResponseController) GetRawData (w http.ResponseWriter, req *http.Request) {
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
	extractor := rc.extractor().Current
	extractorInfo := extractor.Info()

	addBaseHeaders(w.Header())

	bytes, _ := json.Marshal(&extractorInfo)

	_, _ = fmt.Fprint(w, string(bytes))
}
