package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
)

const (
	BinanceWavesBtc      extractors.ExtractorType = "binance-waves-btc"
)

const (
	Endpoint = "https://api.binance.com/api/v3/ticker/price"

	TimeoutSec  = 5
	BtcDecimals = 100000000
)

type Rs struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type Extractor struct{}

func (e *Extractor) Info() *extractors.ExtractorInfo {
	return &extractors.ExtractorInfo{
		Tag:         "binance-WAVES-BTC",
		Description: "Binance extractor",
	}
}
func (e *Extractor) Extract(ctx context.Context) (*extractors.Data, error) {
	endpoint := fmt.Sprintf("%v?symbol=WAVESBTC", Endpoint)

	newCtx, _ := context.WithTimeout(ctx, TimeoutSec*time.Second)
	rq, err := http.NewRequestWithContext(newCtx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	rs, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, err
	}

	defer rs.Body.Close()

	var result Rs

	byteValue, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, err
	}

	v, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return nil, err
	}

	return &extractors.Data{
		Type:  extractors.Int64,
		Value: strconv.FormatInt(int64(v*BtcDecimals), 10),
	}, nil
}
func (e *Extractor) Aggregate(values []extractors.Data) (*extractors.Data, error) {
	var result int64

	for _, data := range values {
		v, err := strconv.ParseInt(data.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		result += v
	}

	return &extractors.Data{
		Type:  extractors.Int64,
		Value: strconv.FormatInt(result/int64(len(values)), 10),
	}, nil
}
