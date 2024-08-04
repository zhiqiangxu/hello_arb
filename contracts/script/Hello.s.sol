// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Script.sol";

contract HelloScript is Script {
    function setUp() public {}

    function run() public pure returns (bytes memory a, bytes memory b) {
        a = abi.encodePacked(uint16(2));
        b = abi.encode(a);
    }
}
