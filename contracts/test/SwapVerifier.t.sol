// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "../src/SwapVerifier.sol";
import "../src/libs/Utils.sol";
import "openzeppelin-contracts/contracts/interfaces/IERC20.sol";

contract SwapVerifierTest is Test {
    uint256 bscFork;
    uint256 fairFork;

    function setUp() public {
        bscFork = vm.createFork("https://bsc-dataseed1.ninicoin.io");
        fairFork = vm.createFork("https://rpc.etherfair.org");
    }

    function testConstructor() public {
        vm.selectFork(bscFork);

        Path[][] memory pathsBatch = new Path[][](1);

        {
            Path[] memory paths = new Path[](2);
            paths[0] = Path({
                exchange: Exchange.UNISWAP,
                pool: vm.parseAddress(
                    "0x0c96e0cacf901f1eb935061caab170f8b2b788dd"
                ),
                to: vm.parseAddress(
                    "0xe9e7cea3dedca5984780bafc599bd69add087d56"
                ),
                fee: 10
            });
            paths[1] = Path({
                exchange: Exchange.UNISWAP,
                pool: vm.parseAddress(
                    "0x66fdb2eccfb58cf098eaa419e5efde841368e489"
                ),
                to: vm.parseAddress(
                    "0x1af3f329e8be154074d8769d1ffa4ee058b1dbc3"
                ),
                fee: 25
            });
            pathsBatch[0] = paths;
        }

        uint256[] memory arbAmountIn = new uint256[](1);
        arbAmountIn[0] = 100;
        uint256[] memory amountOutMin = new uint256[](1);
        amountOutMin[0] = 0;

        address arbToken = 0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3;
        address wtoken = vm.parseAddress(
            "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"
        );
        address router = 0x10ED43C718714eb63d5aA57B78B54704E256024E;
        SwapVerifier v = new SwapVerifier{value: 100}(
            arbToken,
            wtoken,
            router,
            new address[](0),
            arbAmountIn,
            amountOutMin,
            pathsBatch
        );

        checkReturnValue(address(v));
    }

    function checkReturnValue(address v) internal {
        // uint256 size;
        // assembly {
        //     size := extcodesize(v)
        // }
        // bytes memory returnData = new bytes(size);
        // assembly {
        //     extcodecopy(v, add(returnData, 0x20), 0, size)
        // }

        (bool[] memory results, bytes[] memory reasons) = abi.decode(
            address(v).code,
            (bool[], bytes[])
        );

        // console2.log(Utils.bytesToUint256(reasons[0]));

        assertEq(results[0], true, string(reasons[0]));
    }

    function xtestStash() public {
        vm.selectFork(bscFork);

        // // SwapVerifier failed:Pancake: TRANSFER_FAILED
        // bytes
        //     memory args = hex"00000000000000000000000055d398326f99059ff775485246999027b3197955000000000000000000000000bb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c00000000000000000000000010ed43c718714eb63d5aa57b78b54704e256024e00000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000519267a86c84c000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a39af17ce4a8eb807e076805da1e2b8ea7d0755b0000000000000000000000000e09fabb73bd3ade0a17ecc321fd13a19e81ce82000000000000000000000000000000000000000000000000000000000000001900000000000000000000000000000000000000000000000000000000000000000000000000000000000000003764ad558943eb729ebc65cec5996d3e9043af69000000000000000000000000e5434c57e6e1b72f639c728a0889b3d151d31ffa000000000000000000000000000000000000000000000000000000000000001900000000000000000000000000000000000000000000000000000000000000000000000000000000000000007958b8c5f1cf67c9c40b9109a6db8adc3cf471c600000000000000000000000055d398326f99059ff775485246999027b31979550000000000000000000000000000000000000000000000000000000000000019";
        // uint256 value = 0x3fa64a2e810744ca2ea87e;
        // address v = Utils.deployCode(
        //     type(SwapVerifier).creationCode,
        //     args,
        //     value
        // );
        // checkReturnValue(v);

        // SwapVerifier failed:Pancake: OVERFLOW
        vm.prank(vm.parseAddress("c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"));
        bytes
            memory args = hex"000000000000000000000000bb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c000000000000000000000000bb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c00000000000000000000000010ed43c718714eb63d5aa57b78b54704e256024e00000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000280000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000001c000000000000000000000000000000000000000000000000000000000000002e0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000005200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000027fb3ce8c9c7cdca5deeb5e3e486913b97f9d1890000000000000000000000001d2f0da169ceb9fc7b3144628db156f3f6c60dbe000000000000000000000000000000000000000000000000000000000000001e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000005dc30bb8d7f02efef28f7e637d17aea13fa96906000000000000000000000000bb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c00000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003f18135c44c64ebfdcbad8297fe5bdafdbbdd860000000000000000000000001d2f0da169ceb9fc7b3144628db156f3f6c60dbe000000000000000000000000000000000000000000000000000000000000001900000000000000000000000000000000000000000000000000000000000000000000000000000000000000005dc30bb8d7f02efef28f7e637d17aea13fa96906000000000000000000000000bb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c00000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003f18135c44c64ebfdcbad8297fe5bdafdbbdd860000000000000000000000001d2f0da169ceb9fc7b3144628db156f3f6c60dbe0000000000000000000000000000000000000000000000000000000000000019000000000000000000000000000000000000000000000000000000000000000000000000000000000000000027fb3ce8c9c7cdca5deeb5e3e486913b97f9d189000000000000000000000000bb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c000000000000000000000000000000000000000000000000000000000000001e0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000015c921af5f4e89a420f4257cf556c15c333d460d000000000000000000000000a9044f45039b798d21e02cb77f396a2d387cacdd000000000000000000000000000000000000000000000000000000000000001900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000d44ade309e654949b428f8a2f8174041b18a7a0000000000000000000000000bb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000d44ade309e654949b428f8a2f8174041b18a7a0000000000000000000000000a9044f45039b798d21e02cb77f396a2d387cacdd000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000015c921af5f4e89a420f4257cf556c15c333d460d000000000000000000000000bb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c0000000000000000000000000000000000000000000000000000000000000019";
        uint256 value = 0x3fade8d9dd856cd7319f24;
        address v = Utils.deployCode(
            type(SwapVerifier).creationCode,
            args,
            value
        );
        checkReturnValue(v);
    }
}
