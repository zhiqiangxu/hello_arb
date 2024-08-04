// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "../src/Wrapper.sol";
import "../src/libs/Utils.sol";

contract WrapperTest is Test {
    uint256 bscFork;

    function setUp() public {
        bscFork = vm.createFork("https://bsc-dataseed1.ninicoin.io");
    }

    function testEncoding(uint256 amount) public {
        bytes memory amountBytes = abi.encode(amount);

        assertEq(amount, Utils.bytesToUint256(amountBytes));
    }

    function testDeposit() public {
        vm.selectFork(bscFork);

        address wbnb = vm.parseAddress(
            "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"
        );
        uint256 balanceBefore = IERC20(wbnb).balanceOf(address(this));
        uint256 deposit = address(this).balance / 10;
        IWERC20(wbnb).deposit{value: deposit}(deposit);

        uint256 diff = IERC20(wbnb).balanceOf(address(this)) - balanceBefore;

        assertEq(diff, deposit);
    }
}
