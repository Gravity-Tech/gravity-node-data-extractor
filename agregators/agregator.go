package agregators


type Aggregator interface {
	AggregateInt([]interface{}) int64
	AggregateFloat([]interface{}) float64
	AggregateBytes([]interface{}) []byte
	AggregateString([]interface{}) string
}
