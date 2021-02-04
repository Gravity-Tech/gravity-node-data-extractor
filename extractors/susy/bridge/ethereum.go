package bridge

import "github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy"

type EthereumExtractionProvider struct {
	ExtractorDelegate *susy.SourceExtractor
}
