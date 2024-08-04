// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

library ZeroCopySink {
    function WriteBool(bool b) internal pure returns (bytes memory) {
        if (b) {
            return WriteUint8(0x01);
        } else {
            return WriteUint8(0x00);
        }
    }

    function WriteByte(bytes1 b) internal pure returns (bytes memory) {
        return WriteUint8(uint8(b));
    }

    function WriteUint8(uint8 v) internal pure returns (bytes memory) {
        bytes memory buff;

        assembly {
            buff := mload(0x40)
            mstore(buff, 0x01)
            mstore8(add(buff, 0x20), v)
            mstore(0x40, add(buff, 0x21))
        }
        return buff;
    }

    function WriteUint16(uint16 v) internal pure returns (bytes memory) {
        bytes memory buff;

        assembly {
            buff := mload(0x40)
            mstore(buff, 0x02)
            mstore(add(buff, 0x20), shl(240, v))
            mstore(0x40, add(buff, 0x22))
        }

        return buff;
    }

    function WriteUint32(uint32 v) internal pure returns (bytes memory) {
        bytes memory buff;
        assembly {
            buff := mload(0x40)
            mstore(buff, 0x04)
            mstore(add(buff, 0x20), shl(224, v))
            mstore(0x40, add(buff, 0x24))
        }
        return buff;
    }

    function WriteUint64(uint64 v) internal pure returns (bytes memory) {
        bytes memory buff;
        assembly {
            buff := mload(0x40)
            mstore(buff, 0x08)
            mstore(add(buff, 0x20), shl(192, v))
            mstore(0x40, add(buff, 0x28))
        }
        return buff;
    }

    function WriteUint256(uint256 v) internal pure returns (bytes memory) {
        bytes memory buff;

        assembly {
            buff := mload(0x40)
            mstore(buff, 0x20)
            mstore(add(buff, 0x20), v)
            mstore(0x40, add(buff, 0x40))
        }
        return buff;
    }

    function WriteVarUint64(uint64 v) internal pure returns (bytes memory) {
        if (v < 0xFD) {
            return WriteUint8(uint8(v));
        } else if (v <= 0xFFFF) {
            return abi.encodePacked(WriteByte(0xFD), WriteUint16(uint16(v)));
        } else if (v <= 0xFFFFFFFF) {
            return abi.encodePacked(WriteByte(0xFE), WriteUint32(uint32(v)));
        } else {
            return abi.encodePacked(WriteByte(0xFF), WriteUint64(uint64(v)));
        }
    }

    function WriteBytes32(bytes32 data) internal pure returns (bytes memory) {
        bytes memory buff;

        assembly {
            buff := mload(0x40)
            mstore(buff, 0x20)
            mstore(add(buff, 0x20), data)
            mstore(0x40, add(buff, 0x40))
        }

        return buff;
    }

    function WriteBytes20(bytes20 data) internal pure returns (bytes memory) {
        bytes memory buff;

        assembly {
            buff := mload(0x40)
            mstore(buff, 0x14)
            mstore(add(buff, 0x20), data)
            mstore(0x40, add(buff, 0x34))
        }

        return buff;
    }

    function WriteVarBytes(bytes memory data)
        internal
        pure
        returns (bytes memory)
    {
        uint64 l = uint64(data.length);
        return abi.encodePacked(WriteVarUint64(l), data);
    }
}
