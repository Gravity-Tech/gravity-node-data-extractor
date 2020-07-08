package model

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