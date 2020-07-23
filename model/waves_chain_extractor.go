package model

type WavesChainExtractor struct {}

func (e *WavesChainExtractor) DataFeedTag() string {
	return "susy-waves-chain-ext"
}

func (e *WavesChainExtractor) Description() string {
	return "susy waves chain ext"
}

func (e *WavesChainExtractor) Data() (interface{}, interface{}) {
	return nil, nil
}

func (e *WavesChainExtractor) Info() *ExtractorInfo {
	return &ExtractorInfo{
		Description: e.Description(),
		DataFeedTag: e.DataFeedTag(),
	}
}

func (e *WavesChainExtractor) FetchTransferState() *Payment {
	return &Payment{
		ChainName:          "WAVES",
		Amount:             "1000",
		Decimals:           8,
		Timestamp:          10000,
		SourceAddress:      "3b123",
		DestinationAddress: "0x321",
	}
}

func (e *WavesChainExtractor) InitTransfer() {
	//...
}

func (e *WavesChainExtractor) extractData(params interface{}) []RawData {
	return []RawData {}
}

func (e *WavesChainExtractor) mapData(extractedData []RawData) interface{} {
	return nil
}