#/bin/bash -eu

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

pushd $SCRIPT_DIR


GOOS=linux GOARCH=amd64 go build -o arbbot -tags release main.go

mv arbbot release/arbbot
cp -f cmd/artifact/* release/cmd/artifact/

. ~/.bashrc
uploadArbAndCheck

popd