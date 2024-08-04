package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/sethvargo/go-password/password"
	"github.com/urfave/cli"
	"github.com/zhiqiangxu/arbbot/cmd/flag"
)

var MiscCmd = cli.Command{
	Name:  "misc",
	Usage: "misc actions",
	Subcommands: []cli.Command{
		miscBinCmd,
		miscEndianCmd,
		miscPswCmd,
	},
}

var miscBinCmd = cli.Command{
	Name:   "bin",
	Usage:  "show integer in binary form",
	Action: miscBin,
	Flags: []cli.Flag{
		flag.NumbersFlag,
	},
}

var miscEndianCmd = cli.Command{
	Name:   "endian",
	Usage:  "show native endian",
	Action: miscEndian,
	Flags:  []cli.Flag{},
}

var miscPswCmd = cli.Command{
	Name:   "psw",
	Usage:  "generate a password",
	Action: miscPsw,
	Flags:  []cli.Flag{},
}

func miscBin(ctx *cli.Context) (err error) {
	numStrs := strings.Split(ctx.String(flag.NumbersFlag.Name), ",")

	for _, numStr := range numStrs {
		var i int64
		i, err = strconv.ParseInt(numStr, 10, 64)
		if err != nil {
			return
		}
		bin := fmt.Sprintf("%b", uint64(i))
		fmt.Println(bin, "#length", len(bin))
	}
	return
}

func miscEndian(ctx *cli.Context) (err error) {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		fmt.Println("little endian")
	case [2]byte{0xAB, 0xCD}:
		fmt.Println("big endian")
	default:
		panic("Could not determine native endianness.")
	}
	return
}

func miscPsw(ctx *cli.Context) (err error) {
	res, err := password.Generate(20, 10, 10, false, false)
	if err != nil {
		return
	}

	fmt.Println(res)
	return
}
