package abi

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	cabi "github.com/zhiqiangxu/arbbot/contracts/abi/swap_executor"
)

func TestParseFunction(t *testing.T) {

	// reg := regexp.MustCompile(`(?:function\s+)?(\w+)\s*\((.*?)\)\s*(?:returns\s*\((.*?)\))?`)
	// humanReadable := "swapExactTokensForTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline)"
	// matches := reg.FindAllStringSubmatch(humanReadable, -1)

	// match := matches[0]

	// inArgs := strings.TrimSpace(match[2])
	// outArgs := strings.TrimSpace(match[3])
	// reg = regexp.MustCompile(`([^\s,]+)\s+([^\s,]+)`)
	// matches = reg.FindAllStringSubmatch(inArgs, -1)
	// t.Fatal("#matches", len(matches), "inArgs", inArgs, "outArgs", outArgs, "matches", matches)

	// {
	// 	reg := regexp.MustCompile(`(?:event\s+)?(\w+)\s*\((.*?)\)`)
	// 	humanReadable := "event Sync(uint112 reserve0, uint112 reserve1)"
	// 	matches := reg.FindAllStringSubmatch(humanReadable, -1)
	// 	match := matches[0]
	// 	t.Fatal("#match", len(match), match)
	// }

	parsed, _ := abi.JSON(strings.NewReader(cabi.SwapExecutorMetaData.ABI))

	cases := []struct {
		h string
		m string
	}{
		{
			h: "function addOperator(address[] operators)",
			m: "addOperator",
		},
		{
			h: "function balanceOf(address token) returns (uint256 amount)",
			m: "balanceOf",
		},
	}

	for _, c := range cases {
		sig, in, out, err := ParseFunction(c.h)
		if err != nil {
			t.Fatal(err)
		}

		if len(in) != len(parsed.Methods[c.m].Inputs) {
			t.Fatal("#in wrong", in)
		}
		if len(out) != len(parsed.Methods[c.m].Outputs) {
			t.Fatal("#out wrong", out)
		}

		if sig != parsed.Methods[c.m].Sig {
			t.Fatal("sig wrong")
		}
	}

}
