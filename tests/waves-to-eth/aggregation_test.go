package tests

import (
	//waves "github.com/Gravity-Hub-Org/susy-data-extractor/v2/swagger-models/models"

	"bytes"
	"encoding/json"
	controller "github.com/Gravity-Tech/gravity-node-data-extractor/v2/controller"
	_ "go/types"
	"io/ioutil"
	"net/http"
	"testing"
)


type InternalAggregationRequest struct {
	RequestID string
	Receiver string
	Amount int64
}

type InternalAggregationRequestList struct {
	Values []*InternalAggregationRequest
}

func MockupRequestList () *InternalAggregationRequestList {
	return &InternalAggregationRequestList {
		Values: []*InternalAggregationRequest {
			&InternalAggregationRequest{
				RequestID: "A",
				Receiver:  "John",
				Amount:    1000,
			},
			&InternalAggregationRequest{
				RequestID: "A",
				Receiver:  "John",
				Amount:    1500,
			},
			&InternalAggregationRequest{
				RequestID: "A",
				Receiver:  "John",
				Amount:    1000,
			},
			&InternalAggregationRequest{
				RequestID: "A",
				Receiver:  "John",
				Amount:    1000,
			},
		},
	}
}

func MapMockupListToAmounts (values []*InternalAggregationRequest) []int64 {
	result := make([]int64, len(values), len(values))

	for _, value := range values {
		result = append(result, value.Amount)
	}

	return result
}
func MapMockupListAmountToInterfaceList(values []int64) []interface{} {
	mappedAmountList := make([]interface{}, len(values), len(values))

	for i, value := range values {
		mappedAmountList[i] = value
	}
	return mappedAmountList
}


func TestExternalAggregatorResult(t *testing.T) {
	const (
		// THIS IS int64 aggregator
		aggregatorUrl = "https://extractor.gravityhub.org/aggregate"
	)

	mockupValues := MockupRequestList()
	amountList := MapMockupListToAmounts(mockupValues.Values)
	mappedAmountList := MapMockupListAmountToInterfaceList(amountList)

	inputValues := controller.AggregationRequestBody{
		Type:   "int64",
		Values: mappedAmountList,
	}
	requestBody, _ := json.Marshal(&inputValues)

	resp, respErr := http.Post(aggregatorUrl, "application/json", bytes.NewBuffer(requestBody))

	defer resp.Body.Close()

	if respErr != nil {
		t.Errorf("Error. %v \n", respErr)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		t.Errorf("Error. Status: %v; Response Body: %v \n", resp.StatusCode, string(body))
		return
	}

	t.Logf("Aggregated: %v \n", string(body))
}
