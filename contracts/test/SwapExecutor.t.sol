// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "../src/SwapExecutor.sol";
import "../src/libs/Utils.sol";
import "openzeppelin-contracts/contracts/interfaces/IERC20.sol";

contract SwapExecutorTest is Test {
    SwapExecutor public executor;
    uint256 bscFork;

    function setUp() public {
        executor = new SwapExecutor(address(0));

        bscFork = vm.createFork("https://bsc-dataseed1.ninicoin.io");
    }

    function testAddOperator(address[] calldata operators) public {
        executor.addOperator(operators);

        for (uint256 i = 0; i < operators.length; ++i) {
            vm.assume(operators[i] != address(this));
            assertEq(
                executor.hasRole(executor.OPERATOR_ROLE(), operators[i]),
                true,
                "1"
            );
        }

        assertEq(
            executor.hasRole(executor.OPERATOR_ROLE(), address(this)),
            false,
            "2"
        );
        assertEq(executor.owner(), address(this), "3");
    }

    function testRunOutOfGas() public {
        vm.expectRevert();
        vm.prank(address(0));
        payable(0).transfer(address(this).balance);
    }

    function testTargetTx() public {
        vm.selectFork(bscFork);

        address token1 = vm.parseAddress(
            "0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"
        );
        address token0 = vm.parseAddress(
            "0x1d38291068fd6c0e7fc03b77ea90c82950658b37"
        );

        // https://bscscan.com/tx/0xcafe430d8f972977e6b1a47e7f600cd389b29d52c42c8f8d74352077cc10a2c5#eventlog
        address pair = vm.parseAddress(
            "0x6c73bb6841d70310db3e6e33461944505a5da162"
        );
        bytes32 packed = vm.load(pair, bytes32(uint256(8)));
        (, , uint32 unpackedTS) = unpackSlot(uint256(packed));

        uint112 reserve0 = 34934844816753 + 5900000000000;
        uint112 reserve1 = 4770173888281889 - 690693186557990;

        vm.store(
            pair,
            bytes32(uint256(8)),
            bytes32(packSlot(reserve0, reserve1, unpackedTS))
        );

        (uint112 reserve0_, uint112 reserve1_, ) = IPair(pair).getReserves();
        assertEq(reserve0, reserve0_, "1");
        assertEq(reserve1, reserve1_, "2");

        assertEq(IERC20(token0).balanceOf(address(this)), 0, "3");
        deal(token1, pair, reserve1);
        deal(token0, pair, reserve0);

        deal(token1, address(this), 690693186557990);
        IERC20(token1).transfer(pair, 690693186557990);

        IPair(pair).swap(5900000000000, 0, address(this), bytes(""));

        assertEq(IERC20(token0).balanceOf(address(this)), 5900000000000, "4");
    }

    function testDeployAndFlashloan() public {
        vm.selectFork(bscFork);

        SwapExecutor e = new SwapExecutor(address(0));

        address tokenIn = 0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3;
        uint256 amountIn = 100;
        uint256 amountOutMin = 0;

        address otherToken = vm.parseAddress(
            "0xe9e7cea3dedca5984780bafc599bd69add087d56"
        );

        address pool0 = vm.parseAddress(
            "0x0c96e0cacf901f1eb935061caab170f8b2b788dd"
        );
        address pool1 = vm.parseAddress(
            "0x66fdb2eccfb58cf098eaa419e5efde841368e489"
        );
        bool dir = tokenIn < otherToken;

        if (dir) {
            bytes32 packed = vm.load(pool0, bytes32(uint256(8)));
            (, , uint32 unpackedTS) = unpackSlot(uint256(packed));

            vm.store(
                pool0,
                bytes32(uint256(8)),
                bytes32(packSlot(1000, 2000, unpackedTS))
            );

            packed = vm.load(pool1, bytes32(uint256(8)));
            (, , unpackedTS) = unpackSlot(uint256(packed));

            vm.store(
                pool1,
                bytes32(uint256(8)),
                bytes32(packSlot(10000, 2000, unpackedTS))
            );
        } else {
            bytes32 packed = vm.load(pool0, bytes32(uint256(8)));
            (, , uint32 unpackedTS) = unpackSlot(uint256(packed));

            vm.store(
                pool0,
                bytes32(uint256(8)),
                bytes32(packSlot(1000, 2000, unpackedTS))
            );

            packed = vm.load(pool1, bytes32(uint256(8)));
            (, , unpackedTS) = unpackSlot(uint256(packed));

            vm.store(
                pool1,
                bytes32(uint256(8)),
                bytes32(packSlot(1000, 20000, unpackedTS))
            );
        }

        Path[] memory paths = new Path[](2);
        paths[0] = Path({
            exchange: Exchange.UNISWAP,
            pool: pool0,
            to: otherToken,
            fee: 10
        });
        paths[1] = Path({
            exchange: Exchange.UNISWAP,
            pool: pool1,
            to: tokenIn,
            fee: 25
        });

        e.swapFlashloan(tokenIn, amountIn, amountOutMin, paths);
    }

    function testScamToken() public {
        address target = vm.parseAddress(
            // 蜜罐token
            "0x84d76b4f28d44d8a8113a64f2c3af6e57496e585"
        );
        // 正常token
        target = vm.parseAddress("0xe9e7cea3dedca5984780bafc599bd69add087d56");
        runPairWithToken(target);
    }

    function runPairWithToken(address token0) internal {
        vm.selectFork(bscFork);

        address token1 = vm.parseAddress(
            "0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"
        );

        // https://bscscan.com/tx/0xcafe430d8f972977e6b1a47e7f600cd389b29d52c42c8f8d74352077cc10a2c5#eventlog
        address pair = vm.parseAddress(
            "0x6c73bb6841d70310db3e6e33461944505a5da162"
        );
        bytes32 packed = vm.load(pair, bytes32(uint256(8)));
        (, , uint32 unpackedTS) = unpackSlot(uint256(packed));

        uint112 reserve0 = 34934844816753 + 5900000000000;
        uint112 reserve1 = 4770173888281889 - 690693186557990;

        // // modify token0
        vm.store(pair, bytes32(uint256(6)), bytes32(uint256(uint160(token0))));

        vm.store(
            pair,
            bytes32(uint256(8)),
            bytes32(packSlot(reserve0, reserve1, unpackedTS))
        );

        (uint112 reserve0_, uint112 reserve1_, ) = IPair(pair).getReserves();
        assertEq(reserve0, reserve0_, "1");
        assertEq(reserve1, reserve1_, "2");

        assertEq(IERC20(token0).balanceOf(address(this)), 0, "3");

        deal(token1, pair, reserve1);

        deal(token0, pair, reserve0);

        deal(token1, address(this), 690693186557990);
        IERC20(token1).transfer(pair, 690693186557990);

        IPair(pair).swap(5900000000000, 0, address(this), bytes(""));

        assertEq(IERC20(token0).balanceOf(address(this)), 5900000000000, "4");
    }

    function packSlot(
        uint112 reserve0,
        uint112 reserve1,
        uint32 blockTimestampLast
    ) internal pure returns (uint256) {
        return
            uint256(reserve0) +
            (uint256(reserve1) << 112) +
            (uint256(blockTimestampLast) << (112 * 2));
    }

    function unpackSlot(uint256 packed)
        internal
        pure
        returns (
            uint112 unpackedReserve0,
            uint112 unpackedReserve1,
            uint32 unpackedTS
        )
    {
        unpackedReserve0 = uint112(packed);
        unpackedReserve1 = uint112(packed >> 112);
        unpackedTS = uint32(packed >> (112 * 2));
    }

    function testEncode() public {
        uint112 reserve0 = 0x2cc29561b53534941db2f8cdaf1;
        uint112 reserve1 = 0x40d2e942e66f1c1;
        uint32 blockTimestampLast = 0x628f3177;
        uint256 packed = 0x628f3177000000000000040d2e942e66f1c102cc29561b53534941db2f8cdaf1;

        assertEq(packSlot(reserve0, reserve1, blockTimestampLast), packed, "0");

        (
            uint112 unpackedReserve0,
            uint112 unpackedReserve1,
            uint32 unpackedTS
        ) = unpackSlot(packed);

        assertEq(unpackedReserve0, reserve0);
        assertEq(unpackedReserve1, reserve1);
        assertEq(unpackedTS, blockTimestampLast);
    }
}
