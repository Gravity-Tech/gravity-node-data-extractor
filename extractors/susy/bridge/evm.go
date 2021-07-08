package bridge

import "bytes"


func BuildForEVMByteArray(action uint8, rqIDBytes [32]byte,  amountBytes [32]byte, evmReceiver [20]byte) []byte {
	result := []byte{action}
	result = append(result, rqIDBytes[:]...)
	result = append(result, amountBytes[:]...)
	result = append(result, evmReceiver[0:20]...)

	return result
}


func DoesRecordExist(inputAddress []byte) bool {
	if bytes.Equal(inputAddress[:], make([]byte, 32)) || bytes.Equal(inputAddress[:], make([]byte, 20)) {
		return false
	}
	return true
}


type portDelegateClient struct {}

var PortDelegateClient = &portDelegateClient{}

func (pdc *portDelegateClient) PersistByteIX(swapID [32]byte, amount [32]byte, receiver [32]byte, status uint8) []byte {

	var result []byte
	result = append(result, swapID[:]...) 
	result = append(result, amount[:]...) 
	result = append(result, receiver[:]...) 
	result = append(result, status) 

	return result
}
