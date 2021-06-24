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


func ReadSome(vector []byte, offset, len int) []byte {
	return vector[offset:offset + len]
}

func decodeB58Address(addr []byte) [32]byte {
	var resultAddress [32]byte
	copy(resultAddress[:], addr[:])
	return resultAddress
}

func unwrapRequests(encoded []byte, perLength int, swapIds []SwapID, count int) *SwapRequestsDict {
	result := make(SwapRequestsDict)
	i := 0
	
	for i < count {
		encodedRequest := encoded[i*perLength:(i + 1)*perLength]
		decodedRequest, _ := decodeUnwrapRequest(encodedRequest)
		
		fmt.Printf("#%v Foreign: %v \n", 1 + i, decodedRequest.ForeignAddress)
		fmt.Printf("#%v Origin: %v \n", 1 + i, decodedRequest.OriginAddress)
		fmt.Printf("#%v Amount: %v \n", 1 + i, decodedRequest.Amount)
		
		result[swapIds[i]] = decodedRequest
		
		i++
	}
	
	return &result 
}

func unwrapSwapIds(encoded []byte, count int) []SwapID {
	result := make([]SwapID, count)
	i := 0
	
	for i < count {
		var record SwapID
		copy(record[:], encoded[i * 16:(i + 1) * 16])
		result[i] = record
		i++
	}
	
	return result 
}

func decodeUnwrapRequest(encoded []byte) (*IBPortContractUnwrapRequest, int) {
	var internalOffset int

	destination :=  encoded[internalOffset:internalOffset + 32]
	internalOffset += 32
	
	origin := decodeB58Address(encoded[internalOffset:internalOffset + 32])
	internalOffset += 32
		
	amount := binary.LittleEndian.Uint64(encoded[internalOffset:internalOffset + 8])
	
	internalOffset += 8
	
	return &IBPortContractUnwrapRequest{
		OriginAddress:  origin,
		ForeignAddress: destination,
		Amount:         amount,
	}, internalOffset
}

type IBPortContractUnwrapRequest struct {
	OriginAddress      [32]byte
	ForeignAddress     []byte
	Amount             uint64
}

type SwapID [16]byte
type SwapStatusDict map[SwapID]*uint8
type SwapRequestsDict map[SwapID]*IBPortContractUnwrapRequest

// type IBPortContractState struct {
// 	NebulaAddress      [32]byte
// 	TokenAddress       [32]byte
// 	InitializerAddress [32]byte
	
// 	SwapStatusDict     SwapStatusDict
// 	RequestsDict       SwapRequestsDict
// }


type IBPortContractState struct {
	NebulaAddress      solcommon.PublicKey
	TokenAddress       solcommon.PublicKey
	InitializerAddress solcommon.PublicKey
	
	SwapStatusDict     SwapStatusDict
	RequestsDict       SwapRequestsDict
}


func DecodeIBPortState(decoded []byte) *IBPortContractState {
	// decoded, _ := base64.StdEncoding.DecodeString(encodedIBPortState)
	
	// fmt.Println(decoded)
	
	currentOffset := 0
	addressLength := 32
	swapIdLength := 16
	// lengthIndicatorLen := 4
	
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
	
	swapStatusesOffset := 4 + currentOffset + (swapIdLength * int(requestsCount))	
	swapsStatusDict := make(SwapStatusDict)

	var requestIndex uint32

	for requestIndex < requestsCount {
		var swapId SwapID
		
		copy(swapId[:], decoded[currentOffset:currentOffset + swapIdLength])
		
		fmt.Printf("Swap ID: %v \n", swapId)
		
		status := decoded[swapStatusesOffset + int(requestIndex)]
		
		fmt.Printf("Status: %v \n", status)
		statusInt := uint8(status)
		swapsStatusDict[swapId] = &statusInt

		requestIndex++
	}
	
	currentOffset += (int(requestsCount) * swapIdLength) + 4 + (int(requestsCount) * 1)
	
	whatsLeft := decoded[currentOffset:]
	
	swapRequestsCount := int(binary.LittleEndian.Uint32(whatsLeft[0:4]))
	currentOffset += 4

	swapRequestIdsOffset := swapIdLength * swapRequestsCount
	swapRequestIdRanged := decoded[currentOffset:currentOffset + swapRequestIdsOffset]
	currentOffset += swapRequestIdsOffset
	currentOffset += 4
	
	unwrapRequestFlattenedLength := 32 + 32 + 8
	unwrapRequestOffset := unwrapRequestFlattenedLength * swapRequestsCount
	unwrapRequestsRanged := decoded[currentOffset:currentOffset + unwrapRequestOffset]
	
	requestsDict := make(SwapRequestsDict)
	
	decodedSwapIds := unwrapSwapIds(swapRequestIdRanged, swapRequestsCount)
	fmt.Printf("swapRequestIdRanged: %v \n", decodedSwapIds)
	
	fmt.Printf("unwrapRequestsRanged: %v \n", unwrapRequests(unwrapRequestsRanged, unwrapRequestFlattenedLength, decodedSwapIds, swapRequestsCount))
	
	fmt.Printf("whatsLeft: %v \n", whatsLeft)

	fmt.Printf("requestsDict: %v \n", requestsDict)
	fmt.Printf("swapRequestsCount: %v \n", swapRequestsCount)
	
	return nil
}
