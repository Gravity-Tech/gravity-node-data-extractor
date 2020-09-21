package extractors

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

type ExtractorEnumeration = string
type ExtractorEnumerator struct {
	Binance, Metal ExtractorEnumeration
}

var binanceExtractor = "binance"
var metalExtractor = "metal"

var DefaultExtractorEnumerator = &ExtractorEnumerator{
	Binance: binanceExtractor,
	Metal:   metalExtractor,
}

func (e *ExtractorEnumerator) Default() ExtractorEnumeration {
	return binanceExtractor
}

func (e *ExtractorEnumerator) TypeAvailable(enum ExtractorEnumeration) bool {
	available := e.Available()

	for _, item := range available {
		if enum == item { return true }
	}
	return false
}

func (e *ExtractorEnumerator) Available() []ExtractorEnumeration {
	return []string {
		binanceExtractor,
		metalExtractor,
	}
}

type ExtractorProvider struct {
	Current IExtractor
}