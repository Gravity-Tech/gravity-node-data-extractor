package model

// swagger:model
type RawData = byte

type IExtractor interface {
	DataFeedTag() string
	Description() string
	// raw and formated data types
	// first arg should represent type model, second one primitive
	Data() (interface{}, interface{})
	Info() *ExtractorInfo
	extractData(params interface{}) []RawData
	mapData(extractedData []RawData) interface{}
}

// swagger:model
type ExtractorInfo struct {
	Description string `json:"description"`
	DataFeedTag string `json:"datafeedtag"`
}

type ExtractorEnumerator struct {}

var DefaultExtractorEnumerator = &ExtractorEnumerator{}

func (e *ExtractorEnumerator) Available() []string {
	return []string {
		"binance",
		"metal",
	}
}

type ExtractorProvider struct {
	Current IExtractor
}