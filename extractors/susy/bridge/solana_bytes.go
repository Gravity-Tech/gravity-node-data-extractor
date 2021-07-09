package bridge

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/portto/solana-go-sdk/common"
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

type IBPortContractState struct {
	NebulaAddress        solcommon.PublicKey
	TokenAddress         solcommon.PublicKey
	InitializerAddress   solcommon.PublicKey
	Oracles            []solcommon.PublicKey
	
	SwapStatusDict       SwapStatusDict
	RequestsDict         SwapRequestsDict
}

func (swap *SwapID) AsBigInt() *big.Int {
	n := big.NewInt(0)
	n.SetBytes(swap[:])
	return n
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
	
	// fmt.Printf("nebulaAddress: %v \n", base58.Encode(nebulaAddress[:]))

	var tokenAddress [32]byte
	copy(tokenAddress[:], decoded[currentOffset:currentOffset+addressLength])
	currentOffset += addressLength
	
	// fmt.Printf("tokenAddress: %v \n", base58.Encode(tokenAddress[:]))

	var initializerAddress [32]byte
	copy(initializerAddress[:], decoded[currentOffset:currentOffset+addressLength])
	
	currentOffset += addressLength

	oraclesCountBytes := decoded[currentOffset:currentOffset + 4]
	oraclesCount := binary.LittleEndian.Uint32(oraclesCountBytes)

	currentOffset += 4

	// adjustment for new structure
	var oracles []solcommon.PublicKey
	if oraclesCount != 0 {
		currentOffset += (32 * int(oraclesCount))
	} else {
		i := 0
		n := int(oraclesCount)

		for i < n {
			oraclePublicKey := decoded[currentOffset:currentOffset + 32]
			oracles = append(oracles, common.PublicKeyFromBytes(oraclePublicKey))
			currentOffset += 32
			i++
		}
	}

	// fmt.Printf("initializerAddress: %v \n", base58.Encode(initializerAddress[:]))	

	requestsCountBytes := decoded[currentOffset:currentOffset + 4]
	currentOffset += 4
	requestsCount := int(binary.LittleEndian.Uint32(requestsCountBytes))

	fmt.Printf("requestsCount: %v \n", requestsCount)	
	
	swapStatusesOffset := 4 + currentOffset + (swapIdLength * int(requestsCount))	
	swapsStatusDict := make(SwapStatusDict)

	var requestIndex int

	for requestIndex < requestsCount {
		var swapId SwapID
		
		copy(swapId[:], decoded[currentOffset + (requestIndex * swapIdLength):currentOffset + swapIdLength + (requestIndex * swapIdLength)])
		
		// fmt.Printf("Swap ID: %v \n", swapId)
		
		status := decoded[swapStatusesOffset + int(requestIndex)]
		
		// fmt.Printf("Status: %v \n", status)
		statusInt := uint8(status)
		swapsStatusDict[swapId] = &statusInt

		requestIndex++
	}
	
	currentOffset += (int(requestsCount) * swapIdLength) + 4 + (int(requestsCount) * 1)
	
	whatsLeft := decoded[currentOffset:]
	
	swapRequestsCount := int(binary.LittleEndian.Uint32(whatsLeft[0:4]))
	fmt.Printf("swapRequestsCount: %v \n", swapRequestsCount)
	currentOffset += 4

	swapRequestIdsOffset := swapIdLength * swapRequestsCount
	swapRequestIdRanged := decoded[currentOffset:currentOffset + swapRequestIdsOffset]
	currentOffset += swapRequestIdsOffset
	currentOffset += 4
	
	unwrapRequestFlattenedLength := 32 + 32 + 8
	unwrapRequestOffset := unwrapRequestFlattenedLength * swapRequestsCount
	unwrapRequestsRanged := decoded[currentOffset:currentOffset + unwrapRequestOffset]
	
	// requestsDict := make(SwapRequestsDict)
	
	decodedSwapIds := unwrapSwapIds(swapRequestIdRanged, swapRequestsCount)
	// fmt.Printf("swapRequestIdRanged: %v \n", decodedSwapIds)
	
	decodedUnwrapRequests := unwrapRequests(unwrapRequestsRanged, unwrapRequestFlattenedLength, decodedSwapIds, swapRequestsCount)
	// fmt.Printf("unwrapRequestsRanged: %v \n", decodedUnwrapRequests)
	
	// fmt.Printf("whatsLeft: %v \n", whatsLeft)

	// fmt.Printf("requestsDict: %v \n", requestsDict)
	// fmt.Printf("swapRequestsCount: %v \n", swapRequestsCount)
	
	return &IBPortContractState {
		NebulaAddress:      nebulaAddress,
		TokenAddress:       tokenAddress,
		InitializerAddress: initializerAddress,
		Oracles:            oracles,
		
		SwapStatusDict:     swapsStatusDict,
		RequestsDict:       *decodedUnwrapRequests,
	}
}
