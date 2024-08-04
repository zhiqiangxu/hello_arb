#/bin/bash -eu

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

pushd $SCRIPT_DIR

forge inspect SwapExecutor bytecode > abi/artifact/SwapExecutor.bin
forge inspect SwapExecutor abi > abi/artifact/SwapExecutor.abi

forge inspect SwapVerifier bytecode > abi/artifact/SwapVerifier.bin
forge inspect SwapVerifier abi > abi/artifact/SwapVerifier.abi  

forge inspect Wrapper bytecode > abi/artifact/Wrapper.bin
forge inspect Wrapper abi > abi/artifact/Wrapper.abi  

mkdir -p abi/swap_executor
../../go-ethereum/build/bin/abigen --abi abi/artifact/SwapExecutor.abi --bin abi/artifact/SwapExecutor.bin --pkg swap_executor --type SwapExecutor > abi/swap_executor/SwapExecutor.go 

mkdir -p abi/swap_verifier
../../go-ethereum/build/bin/abigen --abi abi/artifact/SwapVerifier.abi --bin abi/artifact/SwapVerifier.bin --pkg swap_verifier --type SwapVerifier > abi/swap_verifier/SwapVerifier.go 

mkdir -p abi/wrapper
../../go-ethereum/build/bin/abigen --abi abi/artifact/Wrapper.abi --bin abi/artifact/Wrapper.bin --pkg wrapper --type Wrapper > abi/wrapper/Wrapper.go 

popd