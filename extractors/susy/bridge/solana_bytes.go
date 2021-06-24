package bridge

import (
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcutil/base58"

	solcommon "github.com/portto/solana-go-sdk/common"
)


type IBPortStateResult struct {
	Data       []string `json:"data"`
	Executable bool     `json:"executable"`
	Lamports   int      `json:"lamports"`
	Owner      string   `json:"owner"`
	RentEpoch  int      `json:"rentEpoch"`
}

type IBPortContractUnwrapRequest struct {
	OriginAddress      [32]byte
	ForeignAddress     [32]byte
	Amount             uint64
}

type SwapID [16]byte
type SwapStatusDict map[SwapID]*uint8
type SwapRequestsDict map[SwapID]*IBPortContractUnwrapRequest

type IBPortContractState struct {
	NebulaAddress      solcommon.PublicKey
	TokenAddress       solcommon.PublicKey
	InitializerAddress solcommon.PublicKey
	
	SwapStatusDict     SwapStatusDict
	RequestsDict       SwapRequestsDict
}

// input string as ba
func DecodeIBPortState(decoded []byte) *IBPortContractState {
	// bodyStr := "i4OiWuJspc63HOthlqpqyqY6N0kh9JKhqI9tDHd5P9EG3fbh12Whk9nL4UbO63msHLSF7V9bN5E6jPWFfv8AqYuDolribKXOtxzrYZaqasqmOjdJIfSSoaiPbQx3eT/RAwAAAMOXxW0T/7WN7R5Bz8HvzCWCmxOASIDk/gmIq6eKjaDxHTseksIWhe9IggV6j09KEAMAAAADAwMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
	// decoded, _ := base64.StdEncoding.DecodeString(bodyStr)
	
	fmt.Println(decoded)
	
	currentOffset := 0
	addressLength := 32
	swapIdLength := 16
	
	var nebulaAddress [32]byte
	copy(nebulaAddress[:], decoded[currentOffset:currentOffset+addressLength])

	currentOffset += addressLength
	
	fmt.Printf("nebulaAddress: %v \n", base58.Encode(nebulaAddress[:]))

	var tokenAddress [32]byte
	copy(tokenAddress[:], decoded[currentOffset:currentOffset+addressLength])
	currentOffset += addressLength
	
	fmt.Printf("tokenAddress: %v \n", base58.Encode(tokenAddress[:]))

	var initializerAddress [32]byte
	copy(initializerAddress[:], decoded[currentOffset:currentOffset+addressLength])
	
	currentOffset += addressLength
	
	fmt.Printf("initializerAddress: %v \n", base58.Encode(initializerAddress[:]))	
			
	requestsCountBytes := decoded[currentOffset:currentOffset + 4]
	currentOffset += 4
	requestsCount := binary.LittleEndian.Uint32(requestsCountBytes)

	fmt.Printf("requestsCount: %v \n", requestsCount)	
	
	var requestIndex uint32

	swapStatusesOffset := 4 + currentOffset + (swapIdLength * int(requestsCount))
	// swapRequestsOffset := 4 + swapStatusesOffset + (swapIdLength * 
	
	swapsStatusDict := make(SwapStatusDict)

	
	for requestIndex < requestsCount {
		var swapId SwapID
		copy(swapId[:], decoded[currentOffset:currentOffset + swapIdLength])
		
		// fmt.Printf("Swap ID: %v \n", swapId)

		status := uint8(decoded[swapStatusesOffset + int(requestIndex)])
		// fmt.Printf("Status: %v \n", status)
		
		swapsStatusDict[swapId] = &status

		currentOffset += swapIdLength
		requestIndex++
	}
	
	fmt.Printf("swapsStatusDict: %v \n", swapsStatusDict)

	return &IBPortContractState{
		NebulaAddress:      nebulaAddress,
		TokenAddress:       tokenAddress,
		InitializerAddress: initializerAddress,
		SwapStatusDict:     swapsStatusDict,
	}
}
