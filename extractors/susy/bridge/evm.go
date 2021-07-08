package bridge


func BuildForEVMByteArray(action uint8, rqIDBytes [32]byte,  amountBytes [32]byte, evmReceiver [20]byte) []byte {
	result := []byte{action}
	result = append(result, rqIDBytes[:]...)
	result = append(result, amountBytes[:]...)
	result = append(result, evmReceiver[0:20]...)

	return result
}