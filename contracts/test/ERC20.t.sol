// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "../src/SwapVerifier.sol";
import "../src/libs/Utils.sol";
import "openzeppelin-contracts/contracts/interfaces/IERC20.sol";

contract ERC20Test is Test {
    uint256 bscFork;
    uint256 ethFork;

    function setUp() public {
        bscFork = vm.createFork("https://bsc-dataseed1.ninicoin.io");
        ethFork = vm.createFork(
            "https://eth-mainnet.nodereal.io/v1/1659dfb40aa24bbb8153a677b98064d7"
        );
    }

    function testSafemoon() public {
        vm.selectFork(bscFork);

        address from = vm.parseAddress(
            "0x0000000000000000000000000000000000000001"
        );
        address to = vm.parseAddress(
            "0x0000000000000000000000000000000000000002"
        );

        // old sfm token
        address token = 0x8076C74C5e3F5852037F31Ff0093Eeb8c8ADd8D3;
        uint256 balanceBefore = IERC20(token).balanceOf(to);
        vm.prank(from);
        // transfer amount is huge
        IERC20(token).transfer(to, 31563346932654954);
        uint256 balanceAfter = IERC20(token).balanceOf(to);

        console2.log(balanceBefore, balanceAfter);
        // but only 2 is received
        console2.log(balanceAfter - balanceBefore);

        assertEq(isFeeOnTransferToken(IERC20(token), from, to, 1000), true);
    }

    function testLeashAndBone() public {
        vm.selectFork(ethFork);

        // leash
        assertEq(
            isFeeOnTransferToken(
                IERC20(
                    vm.parseAddress(
                        "0x27c70cd1946795b66be9d954418546998b546634"
                    )
                ),
                vm.parseAddress("0xa57d319b3cf3ad0e4d19770f71e63cf847263a0b"),
                vm.parseAddress("0x351c3702414c0f1d8acb0156fdc20aeda8d07a8f"),
                1000
            ),
            false
        );

        // bone
        assertEq(
            isFeeOnTransferToken(
                IERC20(
                    vm.parseAddress(
                        "0x9813037ee2218799597d83d4a5b6f3b6778218d9"
                    )
                ),
                vm.parseAddress("0xf7a0383750fef5abace57cc4c9ff98e3790202b3"),
                vm.parseAddress("0xa404f66b9278c4ab8428225014266b4b239bcdc7"),
                1000
            ),
            false
        );
    }

    function isFeeOnTransferToken(
        IERC20 token,
        address from,
        address to,
        uint256 amount
    ) internal returns (bool) {
        uint256 balanceBefore = IERC20(token).balanceOf(to);

        vm.prank(from);
        IERC20(token).transfer(to, amount);

        uint256 balanceAfter = IERC20(token).balanceOf(to);

        return balanceAfter - balanceBefore != amount;
    }
}
