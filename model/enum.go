package model

type ExtractorEnumeration = string
type ExtractorEnumerator struct {
	WavesChain, EthereumChain ExtractorEnumeration
}

var wavesChainExtractor = "waves"
var ethereumChainExtractor = "ethereum"

var DefaultExtractorEnumerator = &ExtractorEnumerator{
	WavesChain: wavesChainExtractor,
	EthereumChain:   ethereumChainExtractor,
}

func (e *ExtractorEnumerator) Default() ExtractorEnumeration {
	return wavesChainExtractor
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
		wavesChainExtractor,
		ethereumChainExtractor,
	}
}

type ExtractorProvider struct {
	Current IExtractor
}