package abi

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var functionRegex = regexp.MustCompile(`(?:function\s+)?(\w+)\s*\((.*?)\)\s*(?:returns\s*\((.*?)\))?`)
var paramsRegex = regexp.MustCompile(`([^\s,]+)\s+([^\s,]+)`)
var eventParamsRegex = regexp.MustCompile(`([^\s,]+)\s+(indexed\s+)?([^\s,]+)`)

var eventRegex = regexp.MustCompile(`(?:event\s+)?(\w+)\s*\((.*?)\)`)

// example humanReadable: swapTokensForExactBNB(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline)
func ParseFunction(humanReadable string) (sig string, in, out abi.Arguments, err error) {

	matches := functionRegex.FindAllStringSubmatch(humanReadable, -1)
	if len(matches) == 0 {
		err = fmt.Errorf("no matches found")
		return
	}
	if len(matches) > 1 {
		err = fmt.Errorf("too many matches found")
		return
	}
	match := matches[0]

	funcName := strings.TrimSpace(match[1])
	inArgs := strings.TrimSpace(match[2])
	var outArgs string
	if len(match) == 4 {
		outArgs = strings.TrimSpace(match[3])
	}

	in, err = parseArgs(inArgs)
	if err != nil {
		return
	}
	out, err = parseArgs(outArgs)
	if err != nil {
		return
	}

	var types []string
	for _, arg := range in {
		types = append(types, arg.Type.String())
	}

	sig = funcName + "(" + strings.Join(types, ",") + ")"
	return
}

func ParseFunctionAsABI(humanReadable string) (ab abi.ABI, err error) {
	sig, in, out, err := ParseFunction(humanReadable)
	if err != nil {
		return
	}

	method := SigToMethod(sig)

	ab.Methods = make(map[string]abi.Method)
	ab.Methods[method] = abi.NewMethod(method, method, abi.Function, "", false, false, in, out)
	return
}

func ParseFunctionsAsABI(humanReadables []string) (ab abi.ABI, err error) {
	ab.Methods = make(map[string]abi.Method)

	for _, humanReadable := range humanReadables {
		var (
			sig     string
			in, out abi.Arguments
		)
		sig, in, out, err = ParseFunction(humanReadable)
		if err != nil {
			return
		}
		method := SigToMethod(sig)

		ab.Methods[method] = abi.NewMethod(method, method, abi.Function, "", false, false, in, out)
	}
	return
}

func ParseEvent(humanReadable string) (sig string, in abi.Arguments, err error) {
	matches := eventRegex.FindAllStringSubmatch(humanReadable, -1)
	if len(matches) == 0 {
		err = fmt.Errorf("no matches found for event")
		return
	}
	if len(matches) > 1 {
		err = fmt.Errorf("too many matches found for event")
		return
	}
	match := matches[0]

	funcName := strings.TrimSpace(match[1])
	inArgs := strings.TrimSpace(match[2])

	in, err = parseEventArgs(inArgs)

	var types []string
	for _, arg := range in {
		types = append(types, arg.Type.String())
	}

	sig = funcName + "(" + strings.Join(types, ",") + ")"
	return
}

func SigToMid(sig string) []byte {
	return crypto.Keccak256([]byte(sig))[0:4]
}

func SigToTopic(sig string) common.Hash {
	return common.BytesToHash(crypto.Keccak256([]byte(sig)))
}

func SigToMethod(sig string) string {
	idx := strings.Index(sig, "(")
	if idx == -1 {
		return sig
	}

	return sig[0:idx]
}

func SigToEvent(sig string) string {
	return SigToMethod(sig)
}

func parseEventArgs(arg string) (args abi.Arguments, err error) {
	if arg == "" {
		return
	}

	matches := eventParamsRegex.FindAllStringSubmatch(arg, -1)
	if len(matches) == 0 {
		return
	}

	for _, match := range matches {
		ty := match[1]
		if ty == "uint" {
			ty = "uint256"
		}
		var abiTy abi.Type
		abiTy, err = abi.NewType(ty, "", nil)
		if err != nil {
			return
		}
		name := match[3]
		args = append(args, abi.Argument{Name: name, Type: abiTy, Indexed: match[2] != ""})
	}
	return
}

func parseArgs(arg string) (args abi.Arguments, err error) {
	if arg == "" {
		return
	}

	matches := paramsRegex.FindAllStringSubmatch(arg, -1)
	if len(matches) == 0 {
		return
	}

	for _, match := range matches {
		ty := match[1]
		if ty == "uint" {
			ty = "uint256"
		}
		var abiTy abi.Type
		abiTy, err = abi.NewType(ty, "", nil)
		if err != nil {
			return
		}
		name := match[2]
		args = append(args, abi.Argument{Name: name, Type: abiTy})
	}
	return
}
