package main

import (
	"os"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli"
	"github.com/zhiqiangxu/arbbot/cmd"
	_ "github.com/zhiqiangxu/arbbot/pkg/metrics"
)

func setupAPP() *cli.App {
	app := cli.NewApp()
	app.Copyright = "Copyright in 2022"
	app.Commands = []cli.Command{
		cmd.ArbCmd,
		cmd.SwapCmd,
		cmd.SwapV3Cmd,
		cmd.Erc20Cmd,
		cmd.EthCmd,
		cmd.SolidityCmd,
		cmd.CompressCmd,
		cmd.FsnCmd,
		cmd.MiscCmd,
	}
	app.Flags = []cli.Flag{}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		zeroLogInit()
		return nil
	}
	return app
}

func consoleWriter() zerolog.ConsoleWriter {
	cw := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339Nano,
	}
	return cw
}

func zeroLogInit() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	log.Logger = zerolog.New(consoleWriter()).With().Timestamp().Logger()
}

func main() {
	if err := setupAPP().Run(os.Args); err != nil {
		log.Fatal().Msg(err.Error())
		os.Exit(1)
	}
}
