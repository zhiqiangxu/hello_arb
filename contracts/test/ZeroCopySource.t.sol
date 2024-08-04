// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "../src/libs/ZeroCopySource.sol";
import "../src/libs/ZeroCopySink.sol";

contract ZeroCopySinkTest is Test {
    function setUp() public {}

    function testNextBool() public {
        bytes memory buff = vm.parseBytes("0x01");
        (bool value, ) = ZeroCopySource.NextBool(buff, 0);
        buff = vm.parseBytes("0x00");
        (value, ) = ZeroCopySource.NextBool(buff, 0);
        assertEq(value, false);

        // uint16 v = 0x0102;
        // assertEq(uint8(v), uint8(0x02));

        // console2.log(vm.toString(bytes2(v)));
        // {
        //     bytes4 buff4 = bytes4(0x01020304);
        //     bytes2 buff2 = bytes2(buff4);
        //     console2.log(vm.toString(buff2));
        // }

        uint256 x = 0x12345678;
        uint8 y = uint8(x);
        console2.log(y);

        uint8 y2;
        assembly {
            let p := mload(0x40)
            mstore(p, x)
            y2 := mload(p)
        }
        console2.log(y2);
    }

    function testNextBoolFuzz(bool v) public {
        bytes memory buff;
        if (v) {
            buff = hex"01";
        } else {
            buff = hex"00";
        }
        (bool value, ) = ZeroCopySource.NextBool(buff, 0);
        assertEq(value, v);
    }

    function testNextByte(bytes memory buff) public {
        vm.assume(buff.length > 0);
        (bytes1 v, ) = ZeroCopySource.NextByte(buff, 0);
        assertEq(v, buff[0]);
    }

    function testNextUint8(bytes memory buff) public {
        vm.assume(buff.length > 0);
        (uint8 v, ) = ZeroCopySource.NextUint8(buff, 0);
        assertEq(v, uint8(buff[0]));
    }

    function testSit(
        bytes1 b1,
        uint8 v8,
        uint16 v16,
        uint32 v32,
        uint64 v64,
        uint64 v256,
        bytes memory bs,
        bytes32 b32,
        bytes20 b20
    ) public {
        {
            (bytes1 decodedB1, ) = ZeroCopySource.NextByte(
                ZeroCopySink.WriteByte(b1),
                0
            );
            assertEq(b1, decodedB1);
        }

        {
            (uint16 decodedV8, ) = ZeroCopySource.NextUint8(
                ZeroCopySink.WriteUint8(v8),
                0
            );
            assertEq(v8, decodedV8);
        }

        {
            (uint16 decodedV16, ) = ZeroCopySource.NextUint16(
                ZeroCopySink.WriteUint16(v16),
                0
            );
            assertEq(v16, decodedV16);
        }

        {
            (uint32 decodedV32, ) = ZeroCopySource.NextUint32(
                ZeroCopySink.WriteUint32(v32),
                0
            );
            assertEq(v32, decodedV32);
        }

        {
            (uint64 decodedV64, ) = ZeroCopySource.NextUint64(
                ZeroCopySink.WriteUint64(v64),
                0
            );
            assertEq(v64, decodedV64);
        }

        {
            (uint256 decodedV256, ) = ZeroCopySource.NextUint256(
                ZeroCopySink.WriteUint256(v256),
                0
            );
            assertEq(v256, decodedV256);
        }

        {
            (uint256 decodedV64, ) = ZeroCopySource.NextVarUint64(
                ZeroCopySink.WriteVarUint64(v64),
                0
            );
            assertEq(v64, decodedV64);
        }

        {
            (bytes memory decodedBs, ) = ZeroCopySource.NextVarBytes(
                ZeroCopySink.WriteVarBytes(bs),
                0
            );
            assertEq(bs, decodedBs);
        }

        {
            (bytes32 decodedB32, ) = ZeroCopySource.NextBytes32(
                ZeroCopySink.WriteBytes32(b32),
                0
            );
            assertEq(b32, decodedB32);
        }

        {
            (bytes20 decodedB20, ) = ZeroCopySource.NextBytes20(
                ZeroCopySink.WriteBytes20(b20),
                0
            );
            assertEq(b20, decodedB20);
        }
    }
}
