package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
	"github.com/zhiqiangxu/arbbot/cmd/flag"
)

var SolidityCmd = cli.Command{
	Name:  "solidity",
	Usage: "solidity actions",
	Subcommands: []cli.Command{
		solidityChangeVersionCmd,
	},
}

var solidityChangeVersionCmd = cli.Command{
	Name:   "change_version",
	Usage:  "change solc version",
	Action: solidityChangeVersion,
	Flags: []cli.Flag{
		flag.PathFlag,
		flag.FromFlag,
		flag.ToFlag,
	},
}

func solidityChangeVersion(ctx *cli.Context) (err error) {
	count := 0
	err = filepath.Walk(ctx.String(flag.PathFlag.Name), func(path string, info fs.FileInfo, _ error) (err error) {
		if info.IsDir() {
			return
		}

		if !strings.HasSuffix(info.Name(), ".sol") {
			return
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return
		}

		content := string(contentBytes)
		if !strings.Contains(content, ctx.String(flag.FromFlag.Name)) {
			return
		}

		fmt.Println(path)

		replaced := strings.Replace(content, ctx.String(flag.FromFlag.Name), ctx.String(flag.ToFlag.Name), 1)

		err = os.WriteFile(path, []byte(replaced), info.Mode())
		if err == nil {
			count++
		}
		return
	})
	if err != nil {
		return
	}

	fmt.Println("#replaced", count)
	return
}
