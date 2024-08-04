package cmd

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"

	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
	"github.com/urfave/cli"
	"github.com/zhiqiangxu/arbbot/cmd/flag"
)

var CompressCmd = cli.Command{
	Name:  "compress",
	Usage: "compress actions",
	Subcommands: []cli.Command{
		gzipEncCmd,
		gzipDecCmd,
		xzEncCmd,
		xzDecCmd,
		zstdEncCmd,
	},
}

var gzipEncCmd = cli.Command{
	Name:   "gz_enc",
	Usage:  "encode with gzip",
	Action: gzipEnc,
	Flags: []cli.Flag{
		flag.InFlag,
		flag.OptionalOutFlag,
	},
}

var zstdEncCmd = cli.Command{
	Name:   "zstd_enc",
	Usage:  "encode with zstd",
	Action: zstdEnc,
	Flags: []cli.Flag{
		flag.InFlag,
		flag.OptionalOutFlag,
	},
}

var gzipDecCmd = cli.Command{
	Name:   "gz_dec",
	Usage:  "decode with gzip",
	Action: gzipDec,
	Flags: []cli.Flag{
		flag.InFlag,
		flag.OptionalOutFlag,
	},
}

var xzEncCmd = cli.Command{
	Name:   "xz_enc",
	Usage:  "encode with xz",
	Action: xzEnc,
	Flags: []cli.Flag{
		flag.InFlag,
		flag.OptionalOutFlag,
	},
}

var xzDecCmd = cli.Command{
	Name:   "xz_dec",
	Usage:  "decode with xz",
	Action: xzDec,
	Flags: []cli.Flag{
		flag.InFlag,
		flag.OptionalOutFlag,
	},
}

func gzipEnc(ctx *cli.Context) (err error) {

	in := ctx.String(flag.InFlag.Name)
	data, err := os.ReadFile(in)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(data)))
	zw, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if err != nil {
		return
	}
	_, err = zw.Write(data)
	if err != nil {
		return
	}
	err = zw.Close()
	if err != nil {
		return
	}

	out := ctx.String(flag.OptionalOutFlag.Name)
	if out == "" {
		out = in
	}

	err = os.WriteFile(out, buf.Bytes(), defaultPermission)
	return
}

func gzipDec(ctx *cli.Context) (err error) {

	in := ctx.String(flag.InFlag.Name)
	data, err := os.ReadFile(in)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(data)
	zr, err := gzip.NewReader(buf)
	if err != nil {
		return
	}
	rbuf := bytes.NewBuffer(make([]byte, 0, len(data)))
	_, err = io.Copy(rbuf, zr)
	if err != nil {
		return
	}

	decoded := rbuf.Bytes()

	out := ctx.String(flag.OptionalOutFlag.Name)
	if out == "" {
		out = in
	}

	err = os.WriteFile(out, decoded, defaultPermission)
	return
}

func xzEnc(ctx *cli.Context) (err error) {
	in := ctx.String(flag.InFlag.Name)
	data, err := os.ReadFile(in)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	w, err := xz.NewWriter(&buf)
	if err != nil {
		return
	}
	if _, err = w.Write(data); err != nil {
		return
	}

	if err = w.Close(); err != nil {
		return
	}

	out := ctx.String(flag.OptionalOutFlag.Name)
	if out == "" {
		out = in
	}

	err = os.WriteFile(out, buf.Bytes(), defaultPermission)
	return

}

func xzDec(ctx *cli.Context) (err error) {
	in := ctx.String(flag.InFlag.Name)
	data, err := os.ReadFile(in)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(data)
	r, err := xz.NewReader(buf)
	if err != nil {
		return
	}
	rbuf := bytes.NewBuffer(make([]byte, 0, len(data)))
	_, err = io.Copy(rbuf, r)
	if err != nil {
		return
	}

	out := ctx.String(flag.OptionalOutFlag.Name)
	if out == "" {
		out = in
	}

	err = os.WriteFile(out, rbuf.Bytes(), defaultPermission)
	return

}

func zstdEnc(ctx *cli.Context) (err error) {
	in := ctx.String(flag.InFlag.Name)
	data, err := os.ReadFile(in)
	if err != nil {
		return
	}

	encoder, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	compressed := encoder.EncodeAll(data, make([]byte, 0, len(data)))

	out := ctx.String(flag.OptionalOutFlag.Name)
	if out == "" {
		out = in
	}

	err = os.WriteFile(out, compressed, defaultPermission)
	return

}

var defaultPermission = os.FileMode(0777)
