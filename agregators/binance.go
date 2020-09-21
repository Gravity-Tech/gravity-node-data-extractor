package agregators

import (
	"fmt"
	"strconv"
)

type BinanceAggregator struct {}

func (aggregator *BinanceAggregator) AggregateInt (values []interface{}) int64 {
	var result int64

	for _, value := range values {
		result += int64(value.(float64))
	}

	return result / int64(len(values))
}

func (aggregator *BinanceAggregator) AggregateFloat (values []interface{}) float64 {
	var result float64

	for _, value := range values {
		result += value.(float64)
	}

	return result / float64(len(values))
}

func (aggregator *BinanceAggregator) AggregateBytes (values []interface{}) []byte {
	return make([]byte, 0)
}

func (aggregator *BinanceAggregator) AggregateString (values []interface{}) string {
	var result float64

	for _, value := range values {
		castedFloat, err := strconv.ParseFloat(value.(string), 64)
		if err != nil { return "" }

		result += castedFloat
	}

	result = result / float64(len(values))

	return string(fmt.Sprintf("%f", result))
}
