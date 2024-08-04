// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "../src/libs/ZeroCopySink.sol";

contract ZeroCopySinkTest is Test {
    function setUp() public {}

    function testWriteUint8(uint8 v) public {
        {
            assertEq(
                vm.toString(ZeroCopySink.WriteUint8(v)),
                vm.toString(abi.encodePacked(v))
            );
        }
    }

    function testWriteUint16(uint16 v) public {
        {
            assertEq(
                vm.toString(ZeroCopySink.WriteUint16(v)),
                vm.toString(abi.encodePacked(v))
            );
        }
    }

    function testWriteUint32(uint32 v) public {
        {
            assertEq(
                vm.toString(ZeroCopySink.WriteUint32(v)),
                vm.toString(abi.encodePacked(v))
            );
        }
    }

    function testWriteUint64(uint64 v) public {
        {
            assertEq(
                vm.toString(ZeroCopySink.WriteUint64(v)),
                vm.toString(abi.encodePacked(v))
            );
        }
    }

    function testWriteUint256(uint256 v) public {
        {
            assertEq(
                vm.toString(ZeroCopySink.WriteUint256(v)),
                vm.toString(abi.encodePacked(v))
            );
        }
    }
}
