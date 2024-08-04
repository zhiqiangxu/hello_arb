package defi

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	boolSliceType   abi.Type
	stringSliceType abi.Type
	arguments       abi.Arguments
)

type SwapVerifierResult struct {
	Success []bool   `json:"success"`
	Reason  [][]byte `json:"reason"`
}

func UnpackVerifierResult(resultBytes []byte, result *SwapVerifierResult) (err error) {
	resultUnpacked, err := arguments.Unpack(resultBytes)
	if err != nil {
		return
	}

	err = arguments.Copy(result, resultUnpacked)
	if err != nil {
		err = fmt.Errorf("arguments.Copy:%v", err)
		return
	}
	return
}

func init() {
	var err error
	boolSliceType, err = abi.NewType("bool[]", "", nil)
	if err != nil {
		panic(err)
	}
	stringSliceType, err = abi.NewType("bytes[]", "", nil)
	if err != nil {
		panic(err)
	}

	arguments = abi.Arguments{
		{Type: boolSliceType, Name: "Success"},
		{Type: stringSliceType, Name: "Reason"},
	}

}

// FYI https://github.com/ethereum/go-ethereum/issues/18360#issuecomment-602180999
// maybe better https://gist.github.com/msigwart/d3e374a64c8718f8ac5ec04b5093597f
var (
	errorSig     = []byte{0x08, 0xc3, 0x79, 0xa0} // Keccak256("Error(string)")[:4]
	abiString, _ = abi.NewType("string", "", nil)
)

func UnpackVerifyError(result []byte) (string, error) {
	if len(result) < 4 || !bytes.Equal(result[:4], errorSig) {
		return string(result), nil
	}
	vs, err := abi.Arguments{{Type: abiString}}.UnpackValues(result[4:])
	if err != nil {
		return "", err
	}
	return vs[0].(string), nil
}
